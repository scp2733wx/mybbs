package file

import "time"

type File struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"index:idx_files_post_id;not null" json:"post_id"`
	UserID    uint      `gorm:"index:idx_files_user_id;not null" json:"user_id"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath  string    `gorm:"type:varchar(255);not null" json:"file_path"`
	Size      int64     `gorm:"not null" json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	FileName  string    `json:"file_name"`
	Size      int64     `json:"size"`
	UpdatedAt time.Time `json:"updated_at"`
}
