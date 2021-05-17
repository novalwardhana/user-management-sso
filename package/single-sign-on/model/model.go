package model

type Result struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type Authorize struct {
	UniqueCode string `json:"unique_code"`
}

type TokenExchangeParams struct {
	Email      string `json:"email"`
	UniqueCode string `json:"unique_code"`
	Domain     string `json:"domain"`
	Secret     string `json:"secret"`
}

type User struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
	IsActive string `json:"is_active" gorm:"column:is_active"`
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
