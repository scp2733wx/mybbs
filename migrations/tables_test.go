package migrations

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// 解析 migrations/tables.sql，提取每张表的列名集合。
// 不连接数据库，纯静态解析，用作 schema vs Go model 一致性回归测试。
func parseTableColumns(sql string) map[string]map[string]bool {
	result := map[string]map[string]bool{}
	re := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+` + "`" + `?(\w+)` + "`" + `?\s*\((.+?)\)\s*ENGINE=`)
	for _, m := range re.FindAllStringSubmatch(sql, -1) {
		table := m[1]
		body := m[2]
		cols := map[string]bool{}
		for _, line := range strings.Split(body, "\n") {
			line = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line), ","))
			if line == "" {
				continue
			}
			upper := strings.ToUpper(line)
			// 跳过非列行：键 / 约束 / 附属 ON 子句
			if strings.HasPrefix(upper, "PRIMARY ") ||
				strings.HasPrefix(upper, "UNIQUE ") ||
				strings.HasPrefix(upper, "KEY ") ||
				strings.HasPrefix(upper, "INDEX ") ||
				strings.HasPrefix(upper, "CONSTRAINT ") ||
				strings.HasPrefix(upper, "FOREIGN ") {
				continue
			}
			f := strings.Fields(line)
			if len(f) == 0 {
				continue
			}
			name := strings.Trim(f[0], "`")
			upperName := strings.ToUpper(name)
			// 继续跳过像 ON UPDATE / COMMENT / COLLATE 这种续行（解析时被拆成一行的情况）
			if name == "" ||
				upperName == "PRIMARY" ||
				upperName == "ON" ||
				upperName == "COMMENT" ||
				upperName == "COLLATE" {
				continue
			}
			cols[name] = true
		}
		result[table] = cols
	}
	return result
}

// 期望列集合，与 api/user/model.go + api/post/model.go 声明的字段一一对应
// （按 GORM snake 命名约定列名；如果以后新增/删除列，必须同时修改这里、SQL、三处 model。）
func expectedColumns() map[string]map[string]bool {
	return map[string]map[string]bool{
		"users": {
			"id": true, "username": true, "name": true, "password_hash": true,
			"role": true, "created_at": true, "updated_at": true, "deleted_at": true,
		},
		"posts": {
			"id": true, "user_id": true, "content": true, "like_count": true,
			"view_count": true, "created_at": true, "updated_at": true, "deleted_at": true,
		},
		"comments": {
			"id": true, "user_id": true, "post_id": true, "content": true,
			"created_at": true, "updated_at": true,
		},
		"post_likes": {
			"user_id": true, "post_id": true, "created_at": true,
		},
	}
}

func TestTablesSQLMatchesModelColumns(t *testing.T) {
	b, err := os.ReadFile("tables.sql")
	if err != nil {
		t.Fatalf("read tables.sql: %v", err)
	}
	actual := parseTableColumns(string(b))
	expected := expectedColumns()

	for table, wantCols := range expected {
		gotCols, ok := actual[table]
		if !ok {
			t.Errorf("table %q missing in tables.sql", table)
			continue
		}
		for col := range wantCols {
			if !gotCols[col] {
				t.Errorf("table %q: expected column %q not found in tables.sql CREATE TABLE", table, col)
			}
		}
		for col := range gotCols {
			if !wantCols[col] {
				t.Errorf("table %q: unexpected column %q present in tables.sql but not declared in expected model set", table, col)
			}
		}
	}
	for table := range actual {
		if _, ok := expected[table]; !ok {
			t.Errorf("tables.sql has table %q not present in expectedColumns mapping", table)
		}
	}
}
