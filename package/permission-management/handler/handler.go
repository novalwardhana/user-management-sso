package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	userVerify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/package/permission-management/model"
	"github.com/novalwardhana/user-management-sso/package/permission-management/usecase"
)

type handler struct {
	usecase usecase.PermissionManagementUsecase
}

func NewHTTPHandler(usecase usecase.PermissionManagementUsecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

func (h *handler) Mount(group *echo.Group) {
	group.GET("/list", h.getPermissionData, userVerify.Verify())
	group.GET("/data/:id", h.getPermissionByID, userVerify.Verify())
	group.POST("/add", h.addPermissionData, userVerify.Verify())
	group.PUT("/update/:id", h.updatePermissionData, userVerify.Verify())
	group.DELETE("/delete/:id", h.deletePermissionData, userVerify.Verify())
}

func (h *handler) getPermissionData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)
	result := <-h.usecase.GetPermissionData()
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Succes get permissions data", Data: result.Data})
}

func (h *handler) getPermissionByID(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	result := <-h.usecase.GetPermissionByID(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success get permission data", Data: result.Data})
}

func (h *handler) addPermissionData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	var newPermission model.NewPermission
	if err := c.Bind(&newPermission); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if newPermission.Code == "" || newPermission.Name == "" || newPermission.Description == "" {
		errMsg := "Code, name, and description must be filled"
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: errMsg})
	}

	newPermission.CreatedAt = time.Now()
	newPermission.UpdatedAt = time.Now()
	result := <-h.usecase.AddPermissionData(newPermission)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success add new permission data", Data: result.Data})
}

func (h *handler) updatePermissionData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	var updatePermission model.UpdatePermission
	if err := c.Bind(&updatePermission); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if updatePermission.Code == "" || updatePermission.Name == "" || updatePermission.Description == "" {
		errMsg := "Code, name, and description must be filled"
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: errMsg})
	}
	updatePermission.UpdatedAt = time.Now()

	result := <-h.usecase.UpdatePermissionData(id, updatePermission)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success update permission data", Data: result.Data})
}

func (h *handler) deletePermissionData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	result := <-h.usecase.DeletePermissionData(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success delete permission data"})
}
