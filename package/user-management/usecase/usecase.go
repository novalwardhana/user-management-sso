package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/novalwardhana/user-management-sso/library/pagination"
	"github.com/novalwardhana/user-management-sso/package/user-management/model"
	"github.com/novalwardhana/user-management-sso/package/user-management/repository"
)

type userManagementUsecase struct {
	repo repository.UserManagementRepo
}

type UserManagementUsecase interface {
	GetUserData(model.ListParams) <-chan model.Result
	GetUserByID(int) <-chan model.Result
	AddUserData(model.NewUserParam) <-chan model.Result
	UpdateUserData(int, model.UpdateUserParam) <-chan model.Result
	DeleteUserData(int) <-chan model.Result
}

func NewUserManagementUsecase(repo repository.UserManagementRepo) UserManagementUsecase {
	return &userManagementUsecase{
		repo: repo,
	}
}

func (uc *userManagementUsecase) GetUserData(params model.ListParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		params.Offset = (params.Page - 1) * params.Limit
		resultTotalData := <-uc.repo.GetTotalUserData(params)
		if resultTotalData.Error != nil {
			output <- model.Result{Error: resultTotalData.Error}
			return
		}

		total := resultTotalData.Data.(int)

		resultData := <-uc.repo.GetUserData(params)
		if resultData.Error != nil {
			output <- model.Result{Error: resultData.Error}
			return
		}
		paginationTable := pagination.PaginationTable{
			Page:        params.Page,
			TotalData:   total,
			DataPerPage: params.Limit,
			Data:        resultData.Data,
		}
		paginationTable.PaginationTotalPage()

		output <- model.Result{Data: paginationTable}
	}()
	return output
}

func (uc *userManagementUsecase) GetUserByID(userID int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.GetUserByID(userID)
		output <- result
	}()
	return output
}

func (uc *userManagementUsecase) AddUserData(param model.NewUserParam) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		param.NewUser.UserUUID = uuid.New().String()
		resultUser := <-uc.repo.AddUserData(param.NewUser)
		if resultUser.Error != nil {
			output <- model.Result{Error: resultUser.Error}
			return
		}

		var userHasRoles []model.UserHasRole
		user := resultUser.Data.(model.NewUser)
		for _, roleID := range param.RoleIDs {
			userHasRole := model.UserHasRole{
				UserID:    user.ID,
				RoleID:    roleID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			userHasRoles = append(userHasRoles, userHasRole)
		}

		resultUserRole := <-uc.repo.AddUserHasRole(userHasRoles)
		if resultUserRole.Error != nil {
			output <- model.Result{Error: resultUserRole.Error}
			return
		}

		output <- model.Result{
			Data: user,
		}
	}()
	return output
}

func (uc *userManagementUsecase) UpdateUserData(id int, param model.UpdateUserParam) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultUser := <-uc.repo.UpdateUserData(id, param.UpdateUser)
		if resultUser.Error != nil {
			output <- model.Result{Error: resultUser.Error}
			return
		}

		resultDeleteRole := <-uc.repo.DeleteUserRoleData(id)
		if resultDeleteRole.Error != nil {
			output <- model.Result{Error: resultDeleteRole.Error}
			return
		}

		var userHasRoles []model.UserHasRole
		for _, roleID := range param.RoleIDs {
			userHasRole := model.UserHasRole{
				UserID:    id,
				RoleID:    roleID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			userHasRoles = append(userHasRoles, userHasRole)
		}

		resultUserRole := <-uc.repo.AddUserHasRole(userHasRoles)
		if resultUserRole.Error != nil {
			output <- model.Result{Error: resultUserRole.Error}
			return
		}

		output <- resultUser
	}()
	return output
}

func (uc *userManagementUsecase) DeleteUserData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultRole := <-uc.repo.DeleteUserRoleData(id)
		if resultRole.Error != nil {
			output <- model.Result{Error: resultRole.Error}
			return
		}

		resultUser := <-uc.repo.DeleteUserData(id)
		output <- resultUser
	}()
	return output
}
