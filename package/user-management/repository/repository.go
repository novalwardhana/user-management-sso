package repository

import (
	"encoding/json"
	"fmt"

	"github.com/novalwardhana/user-management-sso/config/postgres"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/user-management/model"
	log "github.com/sirupsen/logrus"
)

type userManagementRepo struct {
	dbMaster *postgres.DBConnection
}

type UserManagementRepo interface {
	GetUserData() <-chan model.Result
	GetUserByID(int) <-chan model.Result
	AddUserData(model.NewUser) <-chan model.Result
	AddUserHasRole([]model.UserHasRole) <-chan model.Result
	UpdateUserData(int, model.UpdateUser) <-chan model.Result
	DeleteUserData(int) <-chan model.Result
	DeleteUserRoleData(int) <-chan model.Result
}

func NewUserManagementRepo(dbMaster *postgres.DBConnection) UserManagementRepo {
	return &userManagementRepo{
		dbMaster: dbMaster,
	}
}

func (r *userManagementRepo) GetUserData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var listUsers []model.ListUser
		sql := `select 
					u.id,
					u.name,
					u.username,
					u.email,
					u.is_active,
					concat('[', string_agg('{"id":' || r.id ||','|| '"code":"' || r.code ||'","name":"' || r.name || '","group":"' || r."group" || '","description":"' || r.description ||'"}', ','), ']') as roles
				from users u
				left join user_has_roles uhr on u.id = uhr.user_id
				left join roles r on uhr.role_id = r.id
				group by u.id
				order by u.id desc`

		rows, err := r.dbMaster.Read.Raw(sql).Rows()
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		for rows.Next() {
			listUser := model.ListUser{}
			err := rows.Scan(
				&listUser.ID,
				&listUser.Name,
				&listUser.Username,
				&listUser.Email,
				&listUser.IsActive,
				&listUser.Roles,
			)
			if err != nil {
				output <- model.Result{Error: err}
				return
			}
			roles := []map[string]interface{}{}
			_ = json.Unmarshal([]byte(listUser.Roles), &roles)
			listUser.RoleArrays = roles
			listUsers = append(listUsers, listUser)
		}
		output <- model.Result{Data: listUsers}
	}()
	return output
}

func (r *userManagementRepo) GetUserByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var listUser model.ListUser
		sql := `select 
				u.id,
				u.name,
				u.username,
				u.email,
				u.is_active,
				concat('[', string_agg('{"id":' || r.id ||','|| '"code":"' || r.code ||'","name":"' || r.name || '","group":"' || r."group" || '","description":"' || r.description ||'"}', ','), ']') as roles
			from users u
			left join user_has_roles uhr on u.id = uhr.user_id
			left join roles r on uhr.role_id = r.id
			where u.id = ?
			group by u.id order by u.id desc limit 1`
		if err := r.dbMaster.Read.Raw(sql, id).Scan(&listUser).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}

		roles := []map[string]interface{}{}
		if err := json.Unmarshal([]byte(listUser.Roles), &roles); err != nil {
			output <- model.Result{Error: err}
			return
		}
		listUser.RoleArrays = roles
		output <- model.Result{Data: listUser}
	}()
	return output
}

func (r *userManagementRepo) AddUserData(user model.NewUser) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		if err := r.dbMaster.Write.Create(&user).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		user.Password = ""
		user.CreatedAtSTr = user.CreatedAt.Format(constant.DateTimeFormat)
		user.UpdatedAtStr = user.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: user}
	}()
	return output
}

func (r *userManagementRepo) AddUserHasRole(userHasRoles []model.UserHasRole) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		for _, userHasRole := range userHasRoles {
			create := r.dbMaster.Write.Create(&userHasRole)
			if create.Error != nil {
				log.Error(fmt.Sprintf("An error occured when insert user role: %s\n", create.Error.Error()))
			}
		}

		output <- model.Result{}
	}()
	return output
}

func (r *userManagementRepo) UpdateUserData(id int, user model.UpdateUser) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var userTable model.UpdateUser
		var updateData = map[string]interface{}{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"is_active":  user.IsActive,
			"updated_at": user.UpdatedAt,
		}

		update := r.dbMaster.Write.Model(&userTable).Where("id = ?", id).Update(updateData)
		if update.Error != nil {
			output <- model.Result{Error: update.Error}
			return
		}
		if update.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot update, user data not found")}
			return
		}

		user.ID = id
		user.Password = ""
		user.UpdatedAtStr = user.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: user}
	}()
	return output
}

func (r *userManagementRepo) DeleteUserData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var user model.User

		delete := r.dbMaster.Write.Where("id = ?", id).Delete(&user)
		if delete.Error != nil {
			output <- model.Result{Error: delete.Error}
			return
		}
		if delete.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot delete, user data not found")}
			return
		}

		output <- model.Result{}
	}()
	return output
}

func (r *userManagementRepo) DeleteUserRoleData(userID int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var userRole model.UserHasRole
		delete := r.dbMaster.Write.Where("user_id = ?", userID).Delete(&userRole)
		if delete.Error != nil {
			output <- model.Result{Error: delete.Error}
			return
		}

		output <- model.Result{}
	}()
	return output
}
