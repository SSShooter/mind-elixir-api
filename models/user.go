package models

type User struct {
	Id     string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
	From   string `json:"from" binding:"required"`
}

type OriginalUserData interface{}
