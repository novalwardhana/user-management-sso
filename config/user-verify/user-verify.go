package userverify

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/novalwardhana/user-management-sso/global/constant"
	library "github.com/novalwardhana/user-management-sso/library/response"
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

			request := c.Request()
			header := request.Header
			authorization := header.Get("Authorization")
			if len(authorization) <= 0 {
				return c.JSON(http.StatusUnauthorized, library.Response{StatusCode: http.StatusUnauthorized, Message: "Bearer authorization must be filled"})
			}

			authorizationSplit := strings.Split(authorization, " ")
			if len(authorizationSplit) != 2 {
				return c.JSON(http.StatusUnauthorized, library.Response{StatusCode: http.StatusUnauthorized, Message: "Bearer authorization must be filled"})
			}

			token, err := jwt.ParseWithClaims(authorizationSplit[1], &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv(constant.ENVAccessTokenSecret)), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, library.Response{StatusCode: http.StatusUnauthorized, Message: err.Error()})
			}
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, library.Response{StatusCode: http.StatusUnauthorized, Message: "Token is invalid"})
			}

			_, status := token.Claims.(*JwtCustomClaims)
			if !status {
				return c.JSON(http.StatusUnauthorized, library.Response{StatusCode: http.StatusUnauthorized, Message: "Status token is invalid"})
			}

			mc := &RoleContext{Context: c}
			return next(mc)
		}
	}
}
