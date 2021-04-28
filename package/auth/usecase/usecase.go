package usecase

import (
	"fmt"

	"github.com/novalwardhana/user-management-sso/package/auth/model"
	"github.com/novalwardhana/user-management-sso/package/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo repository.AuthRepo
}

type AuthUsecase interface {
	Login(string, string) <-chan model.Result
}

func NewAuthUsecase(repo repository.AuthRepo) AuthUsecase {
	return &authUsecase{
		repo: repo,
	}
}

func (uc *authUsecase) Login(email, password string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultUser := <-uc.repo.GetUserByEmail(email)
		if resultUser.Error != nil {
			output <- model.Result{Error: resultUser.Error}
			return
		}

		user := resultUser.Data.(model.DataUser)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			output <- model.Result{Error: fmt.Errorf("Password not match (%s)", err.Error())}
			return
		}

		output <- model.Result{}

	}()
	return output
}
