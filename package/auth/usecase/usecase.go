package usecase

import (
	"fmt"
	"os"
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
	CreateToken(*model.UserDataToken) (string, error)
	CreateRefreshToken(string) (string, error)
	DecryptRefreshToken(string) <-chan model.Result
	GenerateNewToken(string, string) <-chan model.Result
	GetUserData(int, string) <-chan model.Result
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
		var permissions = resultPermission.Data.(map[string]string)

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
		refreshToken, err := uc.CreateRefreshToken(user.UserUUID)
		if err != nil {
			output <- model.Result{Error: err}
			return
		}

		accessToken := model.AccessToken{
			Type:         "bearer",
			Token:        token,
			RefreshToken: refreshToken,
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

func (uc *authUsecase) CreateToken(data *model.UserDataToken) (string, error) {
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

func (uc *authUsecase) CreateRefreshToken(userUUID string) (string, error) {
	ExpiresIn, err := time.ParseDuration(os.Getenv(constant.ENVRefreshTokenExpiresIn))
	if err != nil {
		ExpiresIn = time.Duration(24 * time.Hour)
	}
	tokenData := userverify.JwtCustomClaims{
		Data: map[string]interface{}{
			"user_uuid": userUUID,
		},
	}
	tokenData.IssuedAt = time.Now().Unix()
	tokenData.ExpiresAt = time.Now().Add(ExpiresIn).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)

	tokenEncrypt, err := token.SignedString([]byte(os.Getenv(constant.ENVRefreshTokenSecret)))
	if err != nil {
		return "", err
	}
	return tokenEncrypt, nil

}

func (uc *authUsecase) DecryptRefreshToken(refreshToken string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		token, err := jwt.ParseWithClaims(refreshToken, &userverify.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv(constant.ENVRefreshTokenSecret)), nil
		})
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		if !token.Valid {
			output <- model.Result{Error: fmt.Errorf("Refresh token is invalid")}
			return
		}

		tokenDecode, status := token.Claims.(*userverify.JwtCustomClaims)
		if !status {
			output <- model.Result{Error: fmt.Errorf("Status token is invalid")}
			return
		}

		tokenData := tokenDecode.Data.(map[string]interface{})
		userUUID := tokenData["user_uuid"]
		if userUUID == nil {
			output <- model.Result{Error: fmt.Errorf("Refresh token is invalid not contain user uuid")}
			return
		}

		output <- model.Result{Data: userUUID.(string)}
	}()
	return output
}

func (uc *authUsecase) GenerateNewToken(email, uuid string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultUser := <-uc.repo.GetUserByEmailUUID(email, uuid)
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
		var roles = resultRole.Data.([]model.Role)
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
		var permissions = resultPermission.Data.(map[string]string)

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
		refreshToken, err := uc.CreateRefreshToken(user.UserUUID)
		if err != nil {
			output <- model.Result{Error: err}
			return
		}

		accessToken := model.AccessToken{
			Type:         "bearer",
			Token:        token,
			RefreshToken: refreshToken,
		}

		output <- model.Result{Data: accessToken}
	}()
	return output
}

func (uc *authUsecase) GetUserData(id int, email string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		resultUser := <-uc.repo.GetUserByEmail(email)
		if resultUser.Error != nil {
			output <- model.Result{Error: fmt.Errorf("Error when get user data: (%s)", resultUser.Error.Error())}
			return
		}
		user := resultUser.Data.(model.User)
		if user.ID != id {
			output <- model.Result{Error: fmt.Errorf("JWT token is invalid")}
			return
		}
		user.Password = ""

		resultRole := <-uc.repo.GetRole(user.ID)
		if resultRole.Error != nil {
			output <- model.Result{Error: fmt.Errorf("Error when get role data: (%s)", resultRole.Error.Error())}
			return
		}
		var roles = resultRole.Data.([]model.Role)
		var roleIDs []int
		var roleMaps = make(map[string]string)
		for _, role := range roles {
			roleIDs = append(roleIDs, role.ID)
			roleMaps[role.Code] = role.Name
		}

		resultPermission := <-uc.repo.GetPermission(roleIDs)
		if resultPermission.Error != nil {
			output <- model.Result{Error: fmt.Errorf("Erro when get permission data: (%s)", resultPermission.Error.Error())}
			return
		}
		var permissions = resultPermission.Data.(map[string]string)

		userDataToken := model.UserDataToken{
			User:        user,
			Roles:       roleMaps,
			Permissions: permissions,
		}

		output <- model.Result{Data: userDataToken}
	}()
	return output
}
