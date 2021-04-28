package model

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

type Response struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

type DataUser struct {
	ID         int                      `json:"id" gorm:"column:id"`
	Name       string                   `json:"name" gorm:"column:name"`
	Username   string                   `json:"username" gorm:"column:username"`
	Email      string                   `json:"email" gorm:"column:email"`
	Password   string                   `json:"password" gorm:"password"`
	IsActive   bool                     `json:"is_active" gorm:"column:is_active"`
	Roles      string                   `json:"-" gorm:"column:roles"`
	RoleArrays []map[string]interface{} `json:"roles" gorm:"-"`
}

func (DataUser) TableName() string {
	return "public.users"
}
