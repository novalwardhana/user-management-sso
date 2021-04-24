package handler

import (
	"crypto/md5"
	"fmt"
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
	group.GET("/users", h.getUserData, userVerify.Verify())
	group.GET("/user/:id", h.getUserByID, userVerify.Verify())
	group.POST("/user/add", h.addUser, userVerify.Verify())
	group.PUT("/user/:id", h.updateUser, userVerify.Verify())
	group.DELETE("/user/:id", h.deleteUser, userVerify.Verify())
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

func (h *handler) addUser(mc echo.Context) error {
	c := mc.(*userVerify.RoleContext)

	newUser := model.NewUser{}
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), constant.PasswordHashCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	newUser.Password = string(encryptedPassword)
	newUser.IsActive = true
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	result := <-h.userManagementUC.AddUserData(newUser)
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
	updateUser := model.UpdateUser{}
	if err := c.Bind(&updateUser); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{StatusCode: http.StatusBadRequest, Message: err.Error()})
	}

	encryptedPassword := fmt.Sprintf("%x", md5.Sum([]byte(updateUser.Password)))
	updateUser.Password = encryptedPassword
	updateUser.UpdatedAt = time.Now()

	result := <-h.userManagementUC.UpdateUserData(id, updateUser)
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
