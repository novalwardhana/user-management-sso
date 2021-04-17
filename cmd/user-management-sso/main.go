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

	userManagementRepo := userManagementRepo.NewUserManagementRepo(dbMasterRead, dbMasterWrite)
	userManagementUsecase := userManagementUsecase.NewUserManagementUsecase(userManagementRepo)
	userManagementHandler := userManagementHandler.NewHTTPHandler(userManagementUsecase)
	userManagementGroup := r.Group("/user-management")
	userManagementHandler.Mount(userManagementGroup)

	port := fmt.Sprintf(":%s", os.Getenv(constant.ENVPort))
	r.Start(port)
}
