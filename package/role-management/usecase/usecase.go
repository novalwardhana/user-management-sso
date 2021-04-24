package usecase

import (
	"github.com/novalwardhana/user-management-sso/package/role-management/model"
	"github.com/novalwardhana/user-management-sso/package/role-management/repository"
)

type roleManagementUsecase struct {
	repo repository.RoleManagementRepo
}

type RoleManagementUsecase interface {
	GetRoleData() <-chan model.Result
	GetRoleByID(int) <-chan model.Result
	AddRoleData(model.NewRole) <-chan model.Result
	UpdateRoleData(int, model.UpdateRole) <-chan model.Result
	DeleteRoleData(int) <-chan model.Result
}

func NewRoleManagementUsecase(repo repository.RoleManagementRepo) RoleManagementUsecase {
	return &roleManagementUsecase{
		repo: repo,
	}
}

func (uc *roleManagementUsecase) GetRoleData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.GetRoleData()
		output <- result
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

func (uc *roleManagementUsecase) AddRoleData(role model.NewRole) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.AddRoleData(role)
		output <- result
	}()
	return output
}

func (uc *roleManagementUsecase) UpdateRoleData(id int, role model.UpdateRole) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.UpdateRoleData(id, role)
		output <- result
	}()
	return output
}

func (uc *roleManagementUsecase) DeleteRoleData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.DeleteRoleData(id)
		output <- result
	}()
	return output
}
