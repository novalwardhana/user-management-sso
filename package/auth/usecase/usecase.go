package usecase

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	userverify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/global/constant"
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
		user := resultUser.Data.(model.User)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			output <- model.Result{Error: fmt.Errorf("Password not match (%s)", err.Error())}
			return
		}
		user.Password = ""

		resultRole := <-uc.repo.GetRole(user.ID)
		if resultRole.Error != nil {
			output <- model.Result{Error: resultRole.Error}
			return
		}
		var roles = resultRole.Data.([]model.Role)

		var roleIDs []int
		for _, role := range roles {
			roleIDs = append(roleIDs, role.ID)
		}

		resultPermission := <-uc.repo.GetPermission(roleIDs)
		if resultPermission.Error != nil {
			output <- model.Result{Error: resultPermission.Error}
			return
		}
		var permissions = resultPermission.Data.([]model.Permission)

		userDataToken := model.UserDataToken{
			User:        user,
			Roles:       roles,
			Permissions: permissions,
		}

		ExpiresIn, err := time.ParseDuration(constant.ENVDefaultTokenExpiresIn)
		if err != nil {
			ExpiresIn = time.Duration(8 * time.Hour)
		}
		tokenData := userverify.JwtCustomClaims{
			Data: userDataToken,
		}
		tokenData.IssuedAt = time.Now().Unix()
		tokenData.ExpiresAt = time.Now().Add(ExpiresIn).Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
		tokenEncrypt, err := token.SignedString([]byte(constant.ENVAccessTokenSecret))
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		accessToken := model.AccessToken{
			Type:  "bearer",
			Token: tokenEncrypt,
		}

		userData := model.UserData{
			User:        user,
			Roles:       roles,
			Permissions: permissions,
			AccessToken: accessToken,
		}

		output <- model.Result{Data: userData}

	}()
	return output
}
