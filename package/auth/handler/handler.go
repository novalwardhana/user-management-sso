package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
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
}

func (h *handler) login(c echo.Context) error {

	email := c.QueryParam("email")
	password := c.QueryParam("password")

	if email == "" || password == "" {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: "Email and password must be filled"})
	}

	result := <-h.usecase.Login(email, password)
	fmt.Println(result)

	return nil
}
