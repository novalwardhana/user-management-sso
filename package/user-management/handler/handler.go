package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/novalwardhana/user-management-sso/package/user-management/usecase"
)

type handler struct {
	userManagementUC usecase.UserManagementUsecase
}

func NewHTTPHandler(usecase usecase.UserManagementUsecase) *handler {
	return &handler{
		userManagementUC: usecase,
	}
}

func (h *handler) Mount(group *echo.Group) {
	group.GET("/users", h.getUserData)
	group.GET("/user/:id", h.getUserByID)
}

func (h *handler) getUserData(c echo.Context) error {
	result := <-h.userManagementUC.GetUserData()
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": result.Error.Error(), "code": http.StatusNotFound})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success", "code": http.StatusOK, "data": result.Data})
}

func (h *handler) getUserByID(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Bad request"})
	}
	result := <-h.userManagementUC.GetUserByID(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Status not found", "code": http.StatusNotFound})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success", "code": http.StatusOK, "data": result.Data})
}
