package userverify

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/novalwardhana/user-management-sso/global/constant"
	library "github.com/novalwardhana/user-management-sso/library/response"
	authModel "github.com/novalwardhana/user-management-sso/package/auth/model"
)

type JwtCustomClaims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

type RoleContext struct {
	User        *authModel.User
	Roles       *map[string]string
	Permissions *map[string]string
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

			tokenDecode, status := token.Claims.(*JwtCustomClaims)
			if !status {
				return c.JSON(http.StatusUnauthorized, library.Response{StatusCode: http.StatusUnauthorized, Message: "Status token is invalid"})
			}

			tokenData := tokenDecode.Data.(map[string]interface{})
			tokenDataUser := tokenData["user"].(map[string]interface{})
			user := authModel.User{
				ID:       int(tokenDataUser["id"].(float64)),
				Name:     tokenDataUser["name"].(string),
				Username: tokenDataUser["username"].(string),
				Email:    tokenDataUser["email"].(string),
				IsActive: tokenDataUser["is_active"].(bool),
			}
			roles := make(map[string]string)
			if tokenData["roles"] != nil {
				tokenDataRoles := tokenData["roles"].(map[string]interface{})
				for index, data := range tokenDataRoles {
					roles[index] = data.(string)
				}
			}
			permissions := make(map[string]string)
			if tokenData["permissions"] != nil {
				tokenDataPermissions := tokenData["permissions"].(map[string]interface{})
				for index, data := range tokenDataPermissions {
					permissions[index] = data.(string)
				}
			}

			mc := &RoleContext{User: &user, Roles: &roles, Permissions: &permissions, Context: c}
			return next(mc)
		}
	}
}
