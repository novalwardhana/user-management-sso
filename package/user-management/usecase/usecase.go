package usecase

import (
	"github.com/novalwardhana/user-management-sso/package/user-management/model"
	"github.com/novalwardhana/user-management-sso/package/user-management/repository"
)

type userManagementUsecase struct {
	repo repository.UserManagementRepo
}

type UserManagementUsecase interface {
	GetUserData() <-chan model.Result
	GetUserByID(int) <-chan model.Result
	AddUserData(model.NewUser) <-chan model.Result
	UpdateUserData(int, model.UpdateUser) <-chan model.Result
	DeleteUserData(int) <-chan model.Result
}

func NewUserManagementUsecase(repo repository.UserManagementRepo) UserManagementUsecase {
	return &userManagementUsecase{
		repo: repo,
	}
}

func (uc *userManagementUsecase) GetUserData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		data := <-uc.repo.GetUserData()
		output <- data
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

func (uc *userManagementUsecase) AddUserData(user model.NewUser) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.AddUserData(user)
		output <- result
	}()
	return output
}

func (uc *userManagementUsecase) UpdateUserData(id int, user model.UpdateUser) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.UpdateUserData(id, user)
		output <- result
	}()
	return output
}

func (uc *userManagementUsecase) DeleteUserData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		result := <-uc.repo.DeleteUserData(id)
		output <- result
	}()
	return output
}
