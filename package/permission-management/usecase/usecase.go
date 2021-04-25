package usecase

import (
	"github.com/novalwardhana/user-management-sso/package/permission-management/model"
	"github.com/novalwardhana/user-management-sso/package/permission-management/repository"
)

type permissionManagementUsecase struct {
	repo repository.PermissionManagementRepo
}

type PermissionManagementUsecase interface {
	GetPermissionData() <-chan model.Result
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

func (uc *permissionManagementUsecase) GetPermissionData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.GetPermissionData()
		output <- result
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
