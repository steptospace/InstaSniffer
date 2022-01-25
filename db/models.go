package db

import (
	"time"
)

type Data struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
