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

type ListUser struct {
	ID         int                      `json:"id" gorm:"column:id"`
	Name       string                   `json:"name" gorm:"column:name"`
	Username   string                   `json:"username" gorm:"column:username"`
	Email      string                   `json:"email" gorm:"column:email"`
	IsActive   bool                     `json:"is_active" gorm:"column:is_active"`
	Roles      string                   `json:"-" gorm:"column:roles"`
	RoleArrays []map[string]interface{} `json:"roles" gorm:"-"`
}

type User struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `json:"email" gorm:"column:email"`
	IsActive bool   `json:"is_active" gorm:"column:is_active"`
}

type NewUserParam struct {
	NewUser
	RoleIDs []int `json:"role_ids" gorm:"-"`
}

type NewUser struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"column:name"`
	Username     string    `json:"username" gorm:"column:username"`
	Email        string    `json:"email" gorm:"column:email"`
	Password     string    `json:"password" gorm:"column:password"`
	IsActive     bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt    time.Time `json:"-" gorm:"column:created_at"`
	CreatedAtSTr string    `json:"created_at" gorm:"-"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtStr string    `json:"updated_at" gorm:"-"`
	UserUUID     string    `json:"-" gorm:"column:user_uuid"`
}

type UserHasRole struct {
	UserID    int       `json:"user_id" gorm:"column:user_id"`
	RoleID    int       `json:"role_id" gorm:"column:role_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type UpdateUserParam struct {
	UpdateUser
	RoleIDs []int `json:"role_ids" gorm:"-"`
}

type UpdateUser struct {
	ID           int       `json:"id" gorm:"-"`
	Name         string    `json:"name" gorm:"column:name"`
	Username     string    `json:"username" gorm:"column:username"`
	Email        string    `json:"email" gorm:"column:email"`
	Password     string    `json:"password" gorm:"column:password"`
	IsActive     bool      `json:"is_active" gorm:"column:is_active"`
	UpdatedAt    time.Time `json:"-" gorm:"column:updated_at"`
	UpdatedAtStr string    `json:"updated_at" gorm:"-"`
}

func (ListUser) TableName() string {
	return "public.users"
}

func (User) TableName() string {
	return "public.users"
}

func (NewUser) TableName() string {
	return "public.users"
}

func (UpdateUser) TableName() string {
	return "public.users"
}

func (UserHasRole) TableName() string {
	return "public.user_has_roles"
}
