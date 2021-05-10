package usecase

import (
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/model"
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/repository"
)

type singleSignOnUsecase struct {
	repo repository.SingleSignOnRepo
}

type SingleSignOnUsecase interface {
	GetUserUUID(int, string) <-chan model.Result
}

func NewSingleSignOnUsecase(repo repository.SingleSignOnRepo) SingleSignOnUsecase {
	return &singleSignOnUsecase{
		repo: repo,
	}
}

func (u *singleSignOnUsecase) GetUserUUID(id int, email string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		result := <-u.repo.GetUserUUID(id, email)
		output <- result
	}()
	return output
}
