package repository

import (
	"github.com/novalwardhana/user-management-sso/config/postgres"
	"github.com/novalwardhana/user-management-sso/package/auth/model"
)

type authRepo struct {
	dbMaster *postgres.DBConnection
}

type AuthRepo interface {
	GetUserByEmail(string) <-chan model.Result
	GetRole(int) <-chan model.Result
	GetPermission([]int) <-chan model.Result
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

		var user model.User
		sql := `select 
				u.id,
				u.name,
				u.username,
				u.email,
				u.password,
				u.is_active,
				u.user_uuid
			from users u
			where u.email = ?
			group by u.id order by u.id desc limit 1`
		err := r.dbMaster.Read.Raw(sql, email).Scan(&user).Error
		if err != nil {
			output <- model.Result{Error: err}
			return
		}

		output <- model.Result{Data: user}
	}()
	return output
}

func (r *authRepo) GetRole(userID int) <-chan model.Result {
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

func (r *authRepo) GetPermission(roleIDs []int) <-chan model.Result {
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
