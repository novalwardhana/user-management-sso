package repository

import (
	"encoding/json"

	"github.com/novalwardhana/user-management-sso/config/postgres"
	"github.com/novalwardhana/user-management-sso/package/auth/model"
)

type authRepo struct {
	dbMaster *postgres.DBConnection
}

type AuthRepo interface {
	GetUserByEmail(string) <-chan model.Result
}

func NewAuthRepo(dbMaster *postgres.DBConnection) AuthRepo {
	return &authRepo{
		dbMaster: dbMaster,
	}
}

func (r *authRepo) GetUserByEmail(email string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var dataUser model.DataUser
		sql := `select 
				u.id,
				u.name,
				u.username,
				u.email,
				u.password,
				u.is_active,
				concat('[', string_agg('{"id":' || r.id ||','|| '"code":"' || r.code ||'","name":"' || r.name || '","group":"' || r."group" || '","description":"' || r.description ||'"}', ','), ']') as roles
			from users u
			left join user_has_roles uhr on u.id = uhr.user_id
			left join roles r on uhr.role_id = r.id
			where u.email = ?
			group by u.id order by u.id desc limit 1`
		err := r.dbMaster.Read.Raw(sql, email).Scan(&dataUser).Error
		if err != nil {
			output <- model.Result{Error: err}
			return
		}

		var roles = []map[string]interface{}{}
		_ = json.Unmarshal([]byte(dataUser.Roles), &roles)
		dataUser.RoleArrays = roles
		output <- model.Result{Data: dataUser}
	}()
	return output
}
