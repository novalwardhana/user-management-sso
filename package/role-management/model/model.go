package model

import "time"

type Result struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type Role struct {
	ID          int    `json:"id" gorm:"column:id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Group       string `json:"group" gorm:"column:group"`
	Description string `json:"description" gorm:"column:description"`
}

type NewRole struct {
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Group        string    `json:"group" gorm:"column:group"`
	Description  string    `json:"description" gorm:"column:description"`
	CreatedAt    time.Time `json:"-" gorm:"column:created_at"`
	CreatedAtStr string    `json:"created_at" gorm:"-"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtStr string    `json:"updated_at" gorm:"-"`
}

type UpdateRole struct {
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Group        string    `json:"group" gorm:"column:group"`
	Description  string    `json:"description" gorm:"column:description"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtSTr string    `json:"updated_at" gorm:"-"`
}

func (NewRole) TableName() string {
	return "public.roles"
}

func (UpdateRole) TableName() string {
	return "public.roles"
}
