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

type Permission struct {
	ID          int    `json:"id" gorm:"column:id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

type NewPermission struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Description  string    `json:"description" gorm:"column:description"`
	CreatedAt    time.Time `json:"-" gorm:"column:created_at"`
	CreatedAtStr string    `json:"created_at" gorm:"-"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtStr string    `json:"updated_at" gorm:"-"`
}

type UpdatePermission struct {
	ID           int       `json:"id" gorm:"-"`
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Description  string    `json:"description" gorm:"column:description"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtStr string    `json:"updated_at" gorm:"-"`
}

func (Permission) TableName() string {
	return "public.permissions"
}

func (NewPermission) TableName() string {
	return "public.permissions"
}

func (UpdatePermission) TableName() string {
	return "public.permissions"
}

type ListParams struct {
	Page   int
	Limit  int
	Offset int
}
