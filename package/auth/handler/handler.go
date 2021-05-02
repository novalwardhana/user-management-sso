package handler

import (
	"net/http"

	"github.com/labstack/echo"
	userverify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/package/auth/model"
	"github.com/novalwardhana/user-management-sso/package/auth/usecase"
)

type handler struct {
	usecase usecase.AuthUsecase
}

func NewHTTPHandler(usecase usecase.AuthUsecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) Mount(group *echo.Group) {
	group.POST("/login", h.login)
	group.POST("/refresh-token", h.refreshToken, userverify.Verify())
}

func (h *handler) login(c echo.Context) error {

	email := c.QueryParam("email")
	password := c.QueryParam("password")

	if email == "" || password == "" {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: "Email and password must be filled"})
	}

	result := <-h.usecase.Login(email, password)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{StatusCode: http.StatusUnauthorized, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success login", Data: result.Data})
}

func (h *handler) refreshToken(mc echo.Context) error {
	return nil
}
