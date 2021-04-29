package userverify

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtCustomClaims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

type RoleContext struct {
	echo.Context
}

func Verify() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mc := &RoleContext{Context: c}
			return next(mc)
		}
	}
}
