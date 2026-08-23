# mybbs

基于 Go (1.26) + Gin + GORM + MySQL 的轻量级校园论坛后端。支持用户（学生/管理员）注册登录、帖子 CRUD、点赞、评论、浏览量与 JWT 鉴权。

---

## 一、目录结构

```
mybbs/
├── api/
│   ├── admin/            管理员能力（强制删帖等）
│   ├── post/             帖子、评论、点赞、详情、列表相关 handler 与 model
│   └── user/             注册、登录、JWT、用户 model
├── cmd/                  main.go（程序入口）
├── config/               config.yaml + viper 加载逻辑
├── database/             MySQL 连接 + AutoMigrate
├── middleware/           JWTAuth / AdminAuth 中间件
├── migrations/           tables.sql（DDL 真源）+ 一致性测试
├── router/               路由注册
├── statehandler/         统一响应 envelope / JWT secret / 启动期致命错误
├── go.mod / go.sum
├── Makefile
└── README.md
```

---

## 二、技术栈

| 层级 | 选择 |
|---|---|
| 语言 | Go 1.26.5 |
| Web 框架 | gin-gonic/gin v1.12 |
| ORM | gorm.io/gorm v1.31 + gorm.io/driver/mysql v1.6 |
| 配置 | spf13/viper v1.21（YAML）|
| 认证 | golang-jwt/jwt v5（HS256）|
| 密码哈希 | golang.org/x/crypto/bcrypt |
| DB | MySQL 8.0（DDL 使用 utf8mb4_0900_ai_ci / DATETIME(3) / ENUM / 外键 / CHECK）|

---

## 三、环境要求

- Go ≥ 1.26
- MySQL ≥ 8.0（8.0.16+ 支持 CHECK 约束；低于该版本 CHECK 会被静默忽略）
- （可选）GNU Make / MinGW make，用来执行 `make build/run/test`；不装就直接用 `go` 原生命令也可

---

## 四、快速开始

### 1. 克隆 & 拉依赖

```bash
git clone <repo> mybbs
cd mybbs
go mod download
```

### 2. 准备数据库

两种建库方式**二选一**，不要混用：

**方式 A：手写 SQL（推荐，结构唯一真源）**

```bash
# 先建库
mysql -uroot -p -e "CREATE DATABASE mybbs DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
# 再跑 DDL
mysql -uroot -p mybbs < migrations/tables.sql
```

**方式 B：仅靠 AutoMigrate**

启动服务时 `cmd/main.go` 会自动 `db.AutoMigrate(...)`，但**不会生成 ENUM、外键、CHECK、自定义索引名**。仅适合本地开发临时跑通。

### 3. 改配置

编辑 `config/config.yaml`：

```yaml
server:
  port: 8080

database:
  enabled: true
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password_here
  name: mybbs
```

> 安全提示：默认 `statehandler.JWTSecret` 是占位符 `"change-me-to-a-random-secret"`。生产部署前必须替换（目前写入代码，建议迁入 `config.yaml: auth.jwt_secret`）。

### 4. 运行

```bash
# 使用 Makefile（在项目根）
make run

# 或者直接 go run（注意：必须在项目根 D:\GO\mybbs 执行，否则相对路径 config/config.yaml 读不到）
go run ./cmd/main.go
```

启动成功日志：
```
[GIN-debug] Listening and serving HTTP on :8080
```

### 5. 健康检查

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"20260001\",\"name\":\"张三\",\"password\":\"Recruit2026\",\"role\":\"student\"}"
```

---

## 五、接口文档

**通用响应结构**（Envelope）：

```json
{ "code": 0, "msg": "Success", "data": { ... } }
```
- `code=0`：成功；
- `code` 非零 = 失败，此时 HTTP 状态码与 `code` 对齐（`SH.Error(c, 401, ...)` 会同时让 HTTP=401 与 body.code=401）；
- 需要登录的接口：在 Header 加 `Authorization: Bearer <access_token>`。

### 公共（无需登录）

| Method | Path | 说明 |
|---|---|---|
| POST | `/api/v1/auth/register` | 注册。Body `{username,name,password,role}`，role 目前被服务端强制为 `student` |
| POST | `/api/v1/auth/login` | 登录。Body `{username,password}`，返回含 `access_token` |

### 登录用户接口（需 JWT）

| Method | Path | 说明 |
|---|---|---|
| POST | `/api/v1/posts` | 发帖。Body `{content}`，content 长度 1-2000 |
| GET  | `/api/v1/posts` | 帖子列表。Query `page=&page_size=&mode=`（mode: `""` 按时间，`hot` 目前占位） |
| GET  | `/api/v1/posts/:post_id` | 帖子详情 + 评论。会将该帖 `view_count += 1` |
| DELETE | `/api/v1/posts/:post_id` | 作者本人删帖（连同评论/点赞） |
| POST | `/api/v1/posts/:post_id/like` | 点赞/取消点赞（点击切换） |
| POST | `/api/v1/posts/likes` | 批量查当前用户对一批帖子是否点赞。Body `{post_ids:[1,2,3]}` |
| POST | `/api/v1/posts/:post_id/comment` | 发评论。Body `{content}` |

### 管理员接口（JWT + role=admin）

| Method | Path | 说明 |
|---|---|---|
| DELETE | `/api/v1/admin/posts/:post_id` | 管理员强制删帖，不限作者 |

---

## 六、常见调用示例（curl / PowerShell）

```powershell
# 注册
$body = @{username='20260001';name='张三';password='Recruit2026'} | ConvertTo-Json
Invoke-RestMethod -Method POST http://127.0.0.1:8080/api/v1/auth/register `
  -ContentType application/json -Body $body

# 登录，保存 token
$loginBody = @{username='20260001';password='Recruit2026'} | ConvertTo-Json
$loginResp = Invoke-RestMethod -Method POST http://127.0.0.1:8080/api/v1/auth/login `
  -ContentType application/json -Body $loginBody
$token = $loginResp.data.data.access_token   # 注意 login.go 目前是双重 envelope，单层后改 $loginResp.data.access_token

# 发帖（带 Bearer token）
$headers = @{ Authorization = "Bearer $token" }
$postBody = @{content='第一次发帖！'} | ConvertTo-Json
Invoke-RestMethod -Method POST http://127.0.0.1:8080/api/v1/posts `
  -Headers $headers -ContentType application/json -Body $postBody

# 看列表
Invoke-RestMethod -Method GET 'http://127.0.0.1:8080/api/v1/posts?page=1&page_size=10' `
  -Headers $headers
```

---

## 七、Makefile 目标

在项目根执行：

| 目标 | 说明 |
|---|---|
| `make build`       | 编译生成 `./bin/mybbs.exe` |
| `make run`         | 直接 `go run ./cmd/main.go`（必须在项目根） |
| `make test`        | `go test ./...`，包括 `migrations/` 下的 SQL vs model 一致性测试 |
| `make vet`         | `go vet ./...` |
| `make lint`        | `gofmt` 检查 + vet |
| `make migrate-sql` | 用 `mysql` 客户端执行 `migrations/tables.sql`（需系统能执行 mysql 命令） |
| `make clean`       | 删除 `bin/` |

---

## 八、测试

目前提供一个纯静态、不连 DB 的回归测试：

```bash
go test -v ./migrations/...
```

`migrations/tables_test.go` 会解析 `migrations/tables.sql` 的 4 张表，与 `api/user/model.go` + `api/post/model.go` 应有列集合做**双向比对**（列不能多也不能少，表不能多也不能少）。防止再出现「model 加了列但 SQL 没加」导致 `Unknown column` 类问题。

---

## 九、常见报错速查

| 日志/现象 | 原因 | 处理 |
|---|---|---|
| `Unknown column 'deleted_at' in 'field list'` | SQL 里缺列但 model 有软删 | 按"四.2 方式 A"重建 `migrations/tables.sql`；或直接 `ALTER TABLE <table> ADD deleted_at DATETIME(3) NULL ADD INDEX idx_deleted_at (deleted_at)` |
| `Unknown column 'like_count' / 'view_count'` | SQL 没加新列 | 同上，或使用最新的 `tables.sql` 重建 |
| `User: unsupported relations for schema Post` | `Post` 里没声明 `User user.User` 关联字段但调用了 `Preload("User")` | 已在 `api/post/model.go` 加了关联；若升级了其他 model 同理 |
| interface conversion nil to uint / panic: `userID.(uint)` | JWT 中间件失败后仍进 handler；或 c.Get 后未 return | 已通过 `SH.Error` 调 `c.Abort()` + 每个错误分支 return 修复，若重现检查是否漏掉 `return` |
| POST register 返回 404 | 路径必须 `POST /api/v1/auth/register`；或 OPTIONS 预检没配 CORS 中间件 | 核对 URL；前后端分离部署时加 `gin-contrib/cors` 中间件 |
| 注册信息获取错误 EOF | Body 空或不是 JSON，或 Content-Type 错 | 按"五、调用示例"填 raw JSON |
| GET `/api/v1/posts` items 每字段都是 0 / 空串 | 已查到 posts，但没把结果映射到 DTO | 见 `api/post/getlist.go`：需补 for 循环赋值 |
| `Headers were already written. Wanted to override status code 401 with 500` | handler 出错前中间件已写过 Response（之前是 panic 遗留） | 当前 `SH.Error` 已 `c.Abort()` 且大多数 return 已补，若出现请检查对应 handler |
| 数据库密码带 `:` 时 DSN 解析失败 | `config.yaml: password: mysql:scp_2733` 这种写法要注意 yaml 冒号解析 | 用引号括起来 `password: "mysql:scp_2733"`；或改成不含特殊字符的密码 |

---

## 十、已知待办 / 技术债

- [ ] `api/post/getlist.go`：`items` 未映射 `posts`，`CommentCount` 没有来源（需 JOIN 或加冗余列），`mode=hot` 目前是 `"nil"`。
- [ ] `api/user/login.go`：成功响应存在双重 envelope、JWT 生成失败后缺 return（上轮已报告）。
- [ ] `statehandler.JWTSecret` 硬编码，需迁入 `config.yaml`。
- [ ] `api/post/comment.go` 写评论后 `posts.comment_count` 未同步维护。
- [ ] `SH.Error` 语义上用于 HTTP 状态，但也承担了业务错误 code 打印；考虑拆成 `BadRequest / Unauthorized / Forbidden / NotFound / Conflict / ServerError` 语义函数。
- [ ] 缺少跨域（CORS）中间件，与浏览器前端联调前必须加入。

---

## 十一、开发约定

- 新增表：先更新 `api/*/model.go` 的 tag → 同步更新 `migrations/tables.sql` → 跑 `go test ./migrations/...`，必须通过。
- 新增 handler：先在 `router.go` 挂路由 → 若需管理员权限就走 `middleware.AdminAuth()` → 错误分支必须 `return`。
- 提交前至少跑过：`go build ./...`、`go vet ./...`、`go test ./...`。
