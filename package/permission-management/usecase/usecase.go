package usecase

import (
	"github.com/novalwardhana/user-management-sso/library/pagination"
	"github.com/novalwardhana/user-management-sso/package/permission-management/model"
	"github.com/novalwardhana/user-management-sso/package/permission-management/repository"
)

type permissionManagementUsecase struct {
	repo repository.PermissionManagementRepo
}

type PermissionManagementUsecase interface {
	GetPermissionData(model.ListParams) <-chan model.Result
	GetPermissionByID(int) <-chan model.Result
	AddPermissionData(model.NewPermission) <-chan model.Result
	UpdatePermissionData(int, model.UpdatePermission) <-chan model.Result
	DeletePermissionData(int) <-chan model.Result
}

func NewPermissionManagementUsecase(repo repository.PermissionManagementRepo) PermissionManagementUsecase {
	return &permissionManagementUsecase{
		repo: repo,
	}
}

func (uc *permissionManagementUsecase) GetPermissionData(params model.ListParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		params.Offset = (params.Page - 1) * params.Limit
		resultTotalData := <-uc.repo.GetTotalPermissionData(params)
		if resultTotalData.Error != nil {
			output <- model.Result{Error: resultTotalData.Error}
			return
		}

		total := resultTotalData.Data.(int)

		resultData := <-uc.repo.GetPermissionData(params)
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

func (uc *permissionManagementUsecase) GetPermissionByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.GetPermissionByID(id)
		output <- result
	}()
	return output
}

func (uc *permissionManagementUsecase) AddPermissionData(permission model.NewPermission) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.AddPermissionData(permission)
		output <- result
	}()
	return output
}

func (uc *permissionManagementUsecase) UpdatePermissionData(id int, permission model.UpdatePermission) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.UpdatePermissionData(id, permission)
		output <- result
	}()
	return output
}

func (uc *permissionManagementUsecase) DeletePermissionData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.DeletePermissionData(id)
		output <- result
	}()
	return output
}
