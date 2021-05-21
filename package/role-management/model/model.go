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

type ListRole struct {
	ID               int                      `json:"id" gorm:"column:id"`
	Code             string                   `json:"code" gorm:"column:code"`
	Name             string                   `json:"name" gorm:"column:name"`
	Group            string                   `json:"group" gorm:"column:group"`
	Description      string                   `json:"description" gorm:"column:description"`
	Permissions      string                   `json:"-" gorm:"column:permissions"`
	PermissionArrays []map[string]interface{} `json:"permissions" gorm:"-"`
}

type Role struct {
	ID          int    `json:"id" gorm:"column:id"`
	Code        string `json:"code" gorm:"column:code"`
	Name        string `json:"name" gorm:"column:name"`
	Group       string `json:"group" gorm:"column:group"`
	Description string `json:"description" gorm:"column:description"`
}

type NewRoleParam struct {
	NewRole
	PermissionIDs []int `json:"permission_ids" gorm:"-"`
}

type NewRole struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Group        string    `json:"group" gorm:"column:group"`
	Description  string    `json:"description" gorm:"column:description"`
	CreatedAt    time.Time `json:"-" gorm:"column:created_at"`
	CreatedAtStr string    `json:"created_at" gorm:"-"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtStr string    `json:"updated_at" gorm:"-"`
}

type RoleHasPermission struct {
	RoleID       int       `json:"role_id" gorm:"column:role_id"`
	PermissionID int       `json:"permission_id" gorm:"column:permission_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type UpdateRoleParam struct {
	UpdateRole
	PermissionIDs []int `json:"permission_ids" gorm:"-"`
}

type UpdateRole struct {
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Group        string    `json:"group" gorm:"column:group"`
	Description  string    `json:"description" gorm:"column:description"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtSTr string    `json:"updated_at" gorm:"-"`
}

func (ListRole) TableName() string {
	return "public.roles"
}

func (NewRole) TableName() string {
	return "public.roles"
}

func (UpdateRole) TableName() string {
	return "public.roles"
}

func (RoleHasPermission) TableName() string {
	return "public.role_has_permissions"
}

type ListParams struct {
	Page   int
	Limit  int
	Offset int
}
