package handler

import (
	"net/http"

	"github.com/labstack/echo"
	userverify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/model"
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/usecase"
)

type handler struct {
	usecase usecase.SingleSignOnUsecase
}

func NewHTTPHandler(usecase usecase.SingleSignOnUsecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) Mount(group *echo.Group) {
	group.GET("/authorize", h.authorize, userverify.Verify())
}

func (h *handler) authorize(mc echo.Context) error {
	c := mc.(*userverify.RoleContext)

	result := <-h.usecase.GetUserUUID(c.User.ID, c.User.Email)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{StatusCode: http.StatusUnauthorized, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Successfully get user unique code", Data: result.Data})
}
