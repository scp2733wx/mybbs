package post

import (
	"time"

	"gorm.io/gorm"

	"mybbs/api/user"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index:idx_comments_user_id;not null" json:"user_id"`
	User      user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
	PostID    uint      `gorm:"index:idx_comments_post_created,sort:desc;not null" json:"post_id"`
	Content   string    `gorm:"type:varchar(1000);not null" json:"content"`
	CreatedAt time.Time `gorm:"index:idx_comments_post_created,sort:desc" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostLike struct {
	UserID    uint      `gorm:"primaryKey;index:idx_post_likes_user_id" json:"user_id"`
	PostID    uint      `gorm:"primaryKey" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index:idx_posts_user_id;not null" json:"user_id"`
	User      user.User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
	Content   string         `gorm:"type:varchar(2000);not null" json:"content"`
	LikeCount int            `gorm:"not null;default:0;check:like_count >= 0" json:"like_count"`
	ViewCount int            `gorm:"not null;default:0;check:view_count >= 0" json:"view_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_posts_deleted_at" json:"-"`
	Comments  []Comment      `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}
