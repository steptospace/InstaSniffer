package db

import (
	"time"
)

type Tasks struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Users struct {
	UId      string `json:"u_id" gorm:"primaryKey"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

type Error struct {
	TaskId      string `json:"task_id"`
	ErrCode     int    `json:"err_code"`
	Description string `json:"description"`
}

type MediaContent struct {
	UserId      string `json:"user_id"`
	IsVideo     bool   `json:"is_video"`
	URL         string `json:"url"`
	Description string `json:"description"`
}
