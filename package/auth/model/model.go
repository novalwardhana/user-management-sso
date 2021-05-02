package model

type Result struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type UserData struct {
	User        User              `json:"user"`
	Roles       map[string]string `json:"roles"`
	Permissions map[string]string `json:"permissions"`
	AccessToken AccessToken       `json:"access_token"`
}

type UserDataToken struct {
	User        User              `json:"user"`
	Roles       map[string]string `json:"roles"`
	Permissions map[string]string `json:"permissions"`
}

type AccessToken struct {
	Type         string `json:"type"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"password"`
	IsActive bool   `json:"is_active" gorm:"column:is_active"`
	UserUUID string `json:"-" gorm:"column:user_uuid"`
}

type Role struct {
	ID    int    `json:"id" gorm:"column:id"`
	Code  string `json:"code" gorm:"column:code"`
	Name  string `json:"name" gorm:"column:name"`
	Group string `json:"group" gorm:"column:group"`
}

type Permission struct {
	ID   int    `json:"id" gorm:"column:id"`
	Code string `json:"code" gorm:"column:code"`
	Name string `json:"name" gorm:"column:name"`
}
