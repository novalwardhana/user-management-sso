package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	constant "github.com/novalwardhana/user-management-sso/global/constant"
	log "github.com/sirupsen/logrus"

	"github.com/novalwardhana/user-management-sso/config/postgres"
	userManagementHandler "github.com/novalwardhana/user-management-sso/package/user-management/handler"
	userManagementRepo "github.com/novalwardhana/user-management-sso/package/user-management/repository"
	userManagementUsecase "github.com/novalwardhana/user-management-sso/package/user-management/usecase"

	roleManagementHandler "github.com/novalwardhana/user-management-sso/package/role-management/handler"
	roleManagementRepo "github.com/novalwardhana/user-management-sso/package/role-management/repository"
	roleManagementUsecase "github.com/novalwardhana/user-management-sso/package/role-management/usecase"

	permissionManagementHandler "github.com/novalwardhana/user-management-sso/package/permission-management/handler"
	permissionManagementRepo "github.com/novalwardhana/user-management-sso/package/permission-management/repository"
	permissionManagementUsecase "github.com/novalwardhana/user-management-sso/package/permission-management/usecase"
)

var (
	dbMasterRead  *gorm.DB
	dbMasterWrite *gorm.DB
)

func main() {
	r := echo.New()
	if err := godotenv.Load("cmd/user-management-sso/.env"); err != nil {
		log.Error(fmt.Sprintf("An error occured: %s", err.Error()))
	}

	dbMasterRead = postgres.DBMasterRead()
	dbMasterWrite = postgres.DBMasterWrite()

	/* User management */
	userManagementRepo := userManagementRepo.NewUserManagementRepo(dbMasterRead, dbMasterWrite)
	userManagementUsecase := userManagementUsecase.NewUserManagementUsecase(userManagementRepo)
	userManagementHandler := userManagementHandler.NewHTTPHandler(userManagementUsecase)
	userManagementGroup := r.Group("/api/v1/user-management")
	userManagementHandler.Mount(userManagementGroup)

	/* Role management */
	roleManagementRepo := roleManagementRepo.NewRoleManagementRepo(dbMasterRead, dbMasterWrite)
	roleManagementUsecase := roleManagementUsecase.NewRoleManagementUsecase(roleManagementRepo)
	roleManagementHandler := roleManagementHandler.NewHTTPHandler(roleManagementUsecase)
	roleManagementGroup := r.Group("/api/v1/role-management")
	roleManagementHandler.Mount(roleManagementGroup)

	/* Permission management */
	permissionManagementRepo := permissionManagementRepo.NewPermissionManagementRepo(dbMasterRead, dbMasterWrite)
	permissionManagementUsecase := permissionManagementUsecase.NewPermissionManagementUsecase(permissionManagementRepo)
	permissionManagementHandler := permissionManagementHandler.NewHTTPHandler(permissionManagementUsecase)
	permissionManagementGroup := r.Group("/api/v1/permission-management")
	permissionManagementHandler.Mount(permissionManagementGroup)

	port := fmt.Sprintf(":%s", os.Getenv(constant.ENVPort))
	r.Start(port)
}
