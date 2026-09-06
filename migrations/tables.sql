DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(32) NOT NULL COMMENT '学号或管理员工号',
    name VARCHAR(32) NOT NULL COMMENT '姓名',
    password_hash VARCHAR(255) NOT NULL COMMENT '安全密码哈希，禁止存储明文',
    role ENUM('student', 'admin') NOT NULL DEFAULT 'student',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '软删除时间，NULL 表示未删除',
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_username (username),
    KEY idx_users_deleted_at (deleted_at),
    CONSTRAINT chk_users_username_not_empty CHECK (CHAR_LENGTH(username) BETWEEN 1 AND 32),
    CONSTRAINT chk_users_name_not_empty CHECK (CHAR_LENGTH(name) BETWEEN 1 AND 32)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE posts (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '帖子作者ID',
    content VARCHAR(2000) NOT NULL COMMENT '帖子正文',
    view_count INT NOT NULL DEFAULT 0 COMMENT '浏览数',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '软删除时间，NULL 表示未删除',
    PRIMARY KEY (id),
    KEY idx_posts_user_id (user_id),
    KEY idx_posts_created_at (created_at DESC, id DESC),
    KEY idx_posts_deleted_at (deleted_at),
    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT chk_posts_content_not_empty CHECK (CHAR_LENGTH(content) BETWEEN 1 AND 2000),
    CONSTRAINT chk_posts_view_count_non_negative CHECK (view_count >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE comments (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '评论作者ID',
    post_id BIGINT UNSIGNED NOT NULL COMMENT '评论帖子ID',
    content VARCHAR(1000) NOT NULL COMMENT '评论内容',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_comments_post_created (post_id, created_at, id),
    KEY idx_comments_user_id (user_id),
    CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id) REFERENCES posts (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT chk_comments_content_not_empty CHECK (CHAR_LENGTH(content) BETWEEN 1 AND 1000)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE likes (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '点赞记录ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '点赞用户ID',
    post_id BIGINT UNSIGNED NOT NULL COMMENT '被点赞帖子ID',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '点赞时间',
    PRIMARY KEY (id),
    UNIQUE KEY idx_user_post_like (user_id, post_id),
    KEY idx_likes_user_id (user_id),
    KEY idx_likes_post_id (post_id),
    CONSTRAINT fk_likes_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT fk_likes_post
        FOREIGN KEY (post_id) REFERENCES posts (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT chk_likes_user_post_not_zero CHECK (user_id > 0 AND post_id > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE files (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    post_id BIGINT UNSIGNED NOT NULL COMMENT '所属帖子ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '上传者用户ID',
    file_name VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_path VARCHAR(255) NOT NULL COMMENT '文件存储路径',
    size BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '上传时间',
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_files_user_id (user_id),
    KEY idx_files_post_id (post_id),
    CONSTRAINT fk_files_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT fk_files_post
        FOREIGN KEY (post_id) REFERENCES posts (id)
        ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT chk_files_file_name_not_empty CHECK (CHAR_LENGTH(file_name) BETWEEN 1 AND 255),
    CONSTRAINT chk_files_file_path_not_empty CHECK (CHAR_LENGTH(file_path) BETWEEN 1 AND 255),
    CONSTRAINT chk_files_size_non_negative CHECK (size >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;