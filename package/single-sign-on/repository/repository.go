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
	GetSSOServiceData(model.TokenExchangeParams) <-chan model.Result
	GerUserByEmailUUID(string, string) <-chan model.Result
	GetRole(int) <-chan model.Result
	GetPermission([]int) <-chan model.Result
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

func (r *singleSignOnRepo) GetSSOServiceData(params model.TokenExchangeParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var amount int
		sql := `select count(*) from sso_services where email = ? and unique_code = ? and domain = ? and secret = ?`
		if err := r.dbMaster.Read.Raw(sql, params.Email, params.UniqueCode, params.Domain, params.Secret).Count(&amount).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		if amount != 1 {
			output <- model.Result{Error: fmt.Errorf("Data not found")}
			return
		}

		output <- model.Result{Data: amount}
	}()
	return output
}

func (r *singleSignOnRepo) GerUserByEmailUUID(email, uuid string) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var user model.User
		sql := `select
				u.id,
				u.name,
				u.username,
				u.email,
				'' as password,
				u.is_active
			from users u
			where u.email = ? and u.user_uuid = ?`
		if err := r.dbMaster.Read.Raw(sql, email, uuid).Scan(&user).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}

		output <- model.Result{Data: user}
	}()
	return output
}

func (r *singleSignOnRepo) GetRole(userID int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var roles []model.Role
		sql := `select 
				r.id,
				r.code,
				r.name,
				r.group
			from user_has_roles uhr 
			inner join roles r on uhr.role_id = r.id
			where uhr.user_id = ?
			group by r.id, r.code, r.name, r.group`
		rows, err := r.dbMaster.Read.Raw(sql, userID).Rows()
		if err != nil {
			output <- model.Result{Error: err}
			return
		}

		for rows.Next() {
			role := model.Role{}
			if err := rows.Scan(
				&role.ID,
				&role.Code,
				&role.Name,
				&role.Group,
			); err != nil {
				output <- model.Result{Error: err}
				return
			}
			roles = append(roles, role)
		}

		output <- model.Result{Data: roles}
	}()
	return output
}

func (r *singleSignOnRepo) GetPermission(roleIDs []int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var permissions = make(map[string]string)
		sql := `select 
				p.id,
				p.code,
				p.name
			from role_has_permissions rhp
			inner join permissions p on rhp.permission_id = p.id 
			where rhp.role_id in (?)
			group by p.id, p.code, p.name`
		rows, err := r.dbMaster.Read.Raw(sql, roleIDs).Rows()
		if err != nil {
			output <- model.Result{Error: err}
			return
		}

		for rows.Next() {
			permission := model.Permission{}
			if err := rows.Scan(
				&permission.ID,
				&permission.Code,
				&permission.Name,
			); err != nil {
				output <- model.Result{Error: err}
				return
			}
			permissions[permission.Code] = permission.Name
		}

		output <- model.Result{Data: permissions}
	}()
	return output
}
