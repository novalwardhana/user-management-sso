package model

type Result struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

type User struct {
	ID       int    `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `json:"email" gorm:"column:email"`
}
