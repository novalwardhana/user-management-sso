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
	group.GET("/token-validation", h.tokenValidation, userverify.Verify())
	group.GET("/token-exchange", h.tokenExchange)
}

func (h *handler) authorize(mc echo.Context) error {
	c := mc.(*userverify.RoleContext)

	result := <-h.usecase.GetUserUUID(c.User.ID, c.User.Email)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{StatusCode: http.StatusUnauthorized, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Successfully get user unique code", Data: result.Data})
}

func (h *handler) tokenValidation(mc echo.Context) error {
	c := mc.(*userverify.RoleContext)
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Successfully validate token"})
}

func (h *handler) tokenExchange(c echo.Context) error {

	email := c.QueryParam("email")
	uniqueCode := c.QueryParam("unique_code")
	domain := c.QueryParam("domain")
	secret := c.QueryParam("secret")

	if len(email) == 0 || len(uniqueCode) == 0 || len(domain) == 0 || len(secret) == 0 {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: "Email, unique code, domain, secret must be filled"})
	}

	params := model.TokenExchangeParams{
		Email:      email,
		UniqueCode: uniqueCode,
		Domain:     domain,
		Secret:     secret,
	}

	result := <-h.usecase.TokenExchange(params)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{StatusCode: http.StatusUnauthorized, Message: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success process token exchange", Data: result.Data})
}
