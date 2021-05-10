package repository

import (
	"fmt"

	"github.com/novalwardhana/user-management-sso/config/postgres"
	"github.com/novalwardhana/user-management-sso/package/single-sign-on/model"
)

type singleSignOnRepo struct {
	dbMaster *postgres.DBConnection
}

type SingleSignOnRepo interface {
	GetUserUUID(int, string) <-chan model.Result
}

func NewSingleSignOnRepo(dbMaster *postgres.DBConnection) SingleSignOnRepo {
	return &singleSignOnRepo{
		dbMaster: dbMaster,
	}
}

func (r *singleSignOnRepo) GetUserUUID(id int, email string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var userUUID []string
		sql := `select user_uuid from users where id = ? and email = ?`

		if err := r.dbMaster.Read.Raw(sql, id, email).Pluck("user_uuid", &userUUID).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}

		if len(userUUID) != 1 {
			output <- model.Result{Error: fmt.Errorf("User UUID not found")}
			return
		}

		var authorize = model.Authorize{
			UniqueCode: userUUID[0],
		}

		output <- model.Result{Data: authorize}
	}()
	return output
}
