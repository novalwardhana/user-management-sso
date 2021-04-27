package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	userVerify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/user-management/model"
	"github.com/novalwardhana/user-management-sso/package/user-management/usecase"
	"golang.org/x/crypto/bcrypt"
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
	group.GET("/list", h.getUserData, userVerify.Verify())
	group.GET("/data/:id", h.getUserByID, userVerify.Verify())
	group.POST("/add", h.addUser, userVerify.Verify())
	group.PUT("/update/:id", h.updateUser, userVerify.Verify())
	group.DELETE("/delete/:id", h.deleteUser, userVerify.Verify())
}

func (h *handler) getUserData(c echo.Context) error {
	result := <-h.userManagementUC.GetUserData()
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success get users data", Data: result.Data})
}

func (h *handler) getUserByID(c echo.Context) error {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	result := <-h.userManagementUC.GetUserByID(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success get user data", Data: result.Data})
}

func (h *handler) addUser(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	var param model.NewUserParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if param.Name == "" || param.Username == "" || param.Email == "" || param.Password == "" {
		errMsg := "Name, username, email, and password must be filled"
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: errMsg})
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), constant.PasswordHashCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	param.Password = string(encryptedPassword)
	param.IsActive = true
	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()

	result := <-h.userManagementUC.AddUserData(param)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success add user data", Data: result.Data})
}

func (h *handler) updateUser(mc echo.Context) error {

	c := mc.(*userVerify.RoleContext)
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	var param model.UpdateUserParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	if param.Name == "" || param.Username == "" || param.Email == "" || param.Password == "" || len(param.RoleIDs) == 0 {
		errMsg := "Name, username, email, password, and role must be filled"
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: errMsg})
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), constant.PasswordHashCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}
	param.Password = string(encryptedPassword)
	param.UpdatedAt = time.Now()

	result := <-h.userManagementUC.UpdateUserData(id, param)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success update user data", Data: result.Data})
}

func (h *handler) deleteUser(mc echo.Context) error {

	c := mc.(*userVerify.RoleContext)
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	result := <-h.userManagementUC.DeleteUserData(id)
	if result.Error != nil {
		return c.JSON(http.StatusNotAcceptable, model.Response{StatusCode: http.StatusNotAcceptable, Message: result.Error.Error()})
	}
	return c.JSON(http.StatusOK, model.Response{StatusCode: http.StatusOK, Message: "Success delete user data"})
}
