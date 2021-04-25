package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	userVerify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/package/role-management/model"
	"github.com/novalwardhana/user-management-sso/package/role-management/usecase"
)

type handler struct {
	usecase usecase.RoleManagementUsecase
}

func NewHTTPHandler(usecase usecase.RoleManagementUsecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) Mount(group *echo.Group) {
	group.GET("/list", h.getRoleData, userVerify.Verify())
	group.GET("/data/:id", h.getRoleByID, userVerify.Verify())
	group.POST("/add", h.addRoleData, userVerify.Verify())
	group.PUT("/update/:id", h.updateRoleData, userVerify.Verify())
	group.DELETE("/delete/:id", h.deleteRoleData, userVerify.Verify())
}

func (h *handler) getRoleData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)
	result := <-h.usecase.GetRoleData()
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success get roles data", Data: result.Data})
}

func (h *handler) getRoleByID(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	result := <-h.usecase.GetRoleByID(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Succes get role data", Data: result.Data})
}

func (h *handler) addRoleData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	var param model.NewRoleParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if param.Code == "" || param.Name == "" || param.Group == "" || param.Description == "" {
		errMsg := "Code, name, group, and description must be filled"
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: errMsg})
	}

	result := <-h.usecase.AddRoleData(param)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success add new role", Data: result.Data})
}

func (h *handler) updateRoleData(mc echo.Context) error {

	c := mc.(*userVerify.RoleContext)

	paramID := c.Param("id")
	var param model.UpdateRoleParam
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if param.Code == "" || param.Name == "" || param.Group == "" || param.Description == "" {
		errMsg := "Code, name, group, and description must be filled"
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: errMsg})
	}
	param.UpdatedAt = time.Now()

	result := <-h.usecase.UpdateRoleData(id, param)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success update role data", Data: result.Data})
}

func (h *handler) deleteRoleData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	result := <-h.usecase.DeleteRoleData(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success delete role data"})
}
