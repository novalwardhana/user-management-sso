package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	userVerify "github.com/novalwardhana/user-management-sso/config/user-verify"
	functionCode "github.com/novalwardhana/user-management-sso/global/function-code"
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
	group.GET("/list", h.getRoleData, userVerify.Verify(functionCode.RoleListData))
	group.GET("/data/:id", h.getRoleByID, userVerify.Verify(functionCode.RoleDetailData))
	group.POST("/add", h.addRoleData, userVerify.Verify(functionCode.RoleAddNewData))
	group.PUT("/update/:id", h.updateRoleData, userVerify.Verify(functionCode.RoleUpdateData))
	group.DELETE("/delete/:id", h.deleteRoleData, userVerify.Verify(functionCode.RoleUpdateData))
}

func (h *handler) getRoleData(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	if page == "" || limit == "" {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: "Page and limit not found"})
	}

	paramPage, err := strconv.Atoi(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	paramLimit, err := strconv.Atoi(limit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	if paramPage <= 0 || paramLimit <= 0 {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: "Page or limit parameter is 0 or less than 0"})
	}

	params := model.ListParams{
		Page:  paramPage,
		Limit: paramLimit,
	}

	result := <-h.usecase.GetRoleData(params)
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
