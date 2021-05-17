package usecase

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	userverify "github.com/novalwardhana/user-management-sso/config/user-verify"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/model"
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/repository"
)

type singleSignOnUsecase struct {
	repo repository.SingleSignOnRepo
}

type SingleSignOnUsecase interface {
	GetUserUUID(int, string) <-chan model.Result
	TokenExchange(model.TokenExchangeParams) <-chan model.Result
	CreateToken(*model.UserDataToken) (string, error)
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

func (uc *singleSignOnUsecase) TokenExchange(params model.TokenExchangeParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultSSOService := <-uc.repo.GetSSOServiceData(params)
		if resultSSOService.Error != nil {
			output <- model.Result{Error: resultSSOService.Error}
			return
		}

		resultUser := <-uc.repo.GerUserByEmailUUID(params.Email, params.UniqueCode)
		if resultUser.Error != nil {
			output <- model.Result{Error: resultUser.Error}
			return
		}
		user := resultUser.Data.(model.User)

		resultRole := <-uc.repo.GetRole(user.ID)
		if resultRole.Error != nil {
			output <- model.Result{Error: resultRole.Error}
			return
		}
		roles := resultRole.Data.([]model.Role)
		var roleIDs []int
		var roleMaps = make(map[string]string)
		for _, role := range roles {
			roleIDs = append(roleIDs, role.ID)
			roleMaps[role.Code] = role.Name
		}

		resultPermission := <-uc.repo.GetPermission(roleIDs)
		if resultPermission.Error != nil {
			output <- model.Result{Error: resultPermission.Error}
			return
		}
		permissions := resultPermission.Data.(map[string]string)

		userDataToken := model.UserDataToken{
			User:        user,
			Roles:       roleMaps,
			Permissions: permissions,
		}
		token, err := uc.CreateToken(&userDataToken)
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		accessToken := model.AccessToken{
			Type:         "bearer",
			Token:        token,
			RefreshToken: "",
		}

		userData := model.UserData{
			User:        user,
			Roles:       roleMaps,
			Permissions: permissions,
			AccessToken: accessToken,
		}

		output <- model.Result{Data: userData}
	}()
	return output
}

func (uc *singleSignOnUsecase) CreateToken(data *model.UserDataToken) (string, error) {
	ExpiresIn, err := time.ParseDuration(os.Getenv(constant.ENVAccessTokenExpiresIn))
	if err != nil {
		ExpiresIn = time.Duration(3 * time.Hour)
	}
	tokenData := userverify.JwtCustomClaims{
		Data: data,
	}
	tokenData.IssuedAt = time.Now().Unix()
	tokenData.ExpiresAt = time.Now().Add(ExpiresIn).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)

	tokenEncrypt, err := token.SignedString([]byte(os.Getenv(constant.ENVAccessTokenSecret)))
	if err != nil {
		return "", err
	}

	return tokenEncrypt, nil
}
