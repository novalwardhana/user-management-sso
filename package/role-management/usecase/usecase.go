package usecase

import (
	"time"

	"github.com/novalwardhana/user-management-sso/library/pagination"
	"github.com/novalwardhana/user-management-sso/package/role-management/model"
	"github.com/novalwardhana/user-management-sso/package/role-management/repository"
)

type roleManagementUsecase struct {
	repo repository.RoleManagementRepo
}

type RoleManagementUsecase interface {
	GetRoleData(model.ListParams) <-chan model.Result
	GetRoleByID(int) <-chan model.Result
	AddRoleData(model.NewRoleParam) <-chan model.Result
	UpdateRoleData(int, model.UpdateRoleParam) <-chan model.Result
	DeleteRoleData(int) <-chan model.Result
}

func NewRoleManagementUsecase(repo repository.RoleManagementRepo) RoleManagementUsecase {
	return &roleManagementUsecase{
		repo: repo,
	}
}

func (uc *roleManagementUsecase) GetRoleData(params model.ListParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		params.Offset = (params.Page - 1) * params.Limit
		resultTotalData := <-uc.repo.GetTotalRoleData(params)
		if resultTotalData.Error != nil {
			output <- model.Result{Error: resultTotalData.Error}
			return
		}
		total := resultTotalData.Data.(int)

		resultData := <-uc.repo.GetRoleData(params)
		if resultData.Error != nil {
			output <- model.Result{Error: resultTotalData.Error}
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

func (uc *roleManagementUsecase) GetRoleByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.GetRoleByID(id)
		output <- result
	}()
	return output
}

func (uc *roleManagementUsecase) AddRoleData(param model.NewRoleParam) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultRole := <-uc.repo.AddRoleData(param.NewRole)
		if resultRole.Error != nil {
			output <- model.Result{Error: resultRole.Error}
			return
		}

		var roleHasPermissions []model.RoleHasPermission
		role := resultRole.Data.(model.NewRole)
		for _, permissionID := range param.PermissionIDs {
			roleHasPermission := model.RoleHasPermission{
				RoleID:       role.ID,
				PermissionID: permissionID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			roleHasPermissions = append(roleHasPermissions, roleHasPermission)
		}

		resultRolePermission := <-uc.repo.AddRolePermissionData(roleHasPermissions)
		if resultRolePermission.Error != nil {
			output <- model.Result{Error: resultRolePermission.Error}
			return
		}

		output <- model.Result{
			Data: role,
		}
	}()
	return output
}

func (uc *roleManagementUsecase) UpdateRoleData(id int, param model.UpdateRoleParam) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultRole := <-uc.repo.UpdateRoleData(id, param.UpdateRole)
		if resultRole.Error != nil {
			output <- model.Result{Error: resultRole.Error}
			return
		}

		resultDeletePermission := <-uc.repo.DeleteRolePermissionData(id)
		if resultDeletePermission.Error != nil {
			output <- model.Result{Error: resultDeletePermission.Error}
			return
		}

		var roleHasPermissions []model.RoleHasPermission
		for _, permissionID := range param.PermissionIDs {
			roleHasPermission := model.RoleHasPermission{
				RoleID:       id,
				PermissionID: permissionID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			roleHasPermissions = append(roleHasPermissions, roleHasPermission)
		}
		resultRolePermission := <-uc.repo.AddRolePermissionData(roleHasPermissions)
		if resultRolePermission.Error != nil {
			output <- model.Result{Error: resultRolePermission.Error}
			return
		}

		output <- resultRole
	}()
	return output
}

func (uc *roleManagementUsecase) DeleteRoleData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultPermission := <-uc.repo.DeleteRolePermissionData(id)
		if resultPermission.Error != nil {
			output <- model.Result{Error: resultPermission.Error}
			return
		}

		resultRole := <-uc.repo.DeleteRoleData(id)
		output <- resultRole
	}()
	return output
}
