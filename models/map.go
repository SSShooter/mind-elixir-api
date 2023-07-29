package models

import "time"

type Map struct {
	Name     string                 `json:"name" binding:"required"`
	Content  map[string]interface{} `json:"content" binding:"required"`
	Author   string                 `json:"author" binding:"required"`
	Date     time.Time              `json:"date"`
	UpdateAt time.Time              `json:"updateAt"`
}
