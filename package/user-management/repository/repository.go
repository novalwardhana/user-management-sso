package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/user-management/model"
)

type userManagementRepo struct {
	dbMasterRead  *gorm.DB
	dbMasterWrite *gorm.DB
}

type UserManagementRepo interface {
	GetUserData() <-chan model.Result
	GetUserByID(int) <-chan model.Result
	AddUserData(model.NewUser) <-chan model.Result
	UpdateUserData(int, model.UpdateUser) <-chan model.Result
	DeleteUserData(int) <-chan model.Result
}

func NewUserManagementRepo(dbMasterRead, dbMasterWrite *gorm.DB) UserManagementRepo {
	return &userManagementRepo{
		dbMasterRead:  dbMasterRead,
		dbMasterWrite: dbMasterWrite,
	}
}

func (r *userManagementRepo) GetUserData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var users []model.User
		sql := "SELECT id, name, username, email, is_active FROM users"
		rows, err := r.dbMasterRead.Raw(sql).Rows()
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		for rows.Next() {
			user := model.User{}
			err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.Username,
				&user.Email,
				&user.IsActive,
			)
			if err != nil {
				output <- model.Result{Error: err}
				return
			}
			users = append(users, user)
		}
		output <- model.Result{Data: users}
	}()
	return output
}

func (r *userManagementRepo) GetUserByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var user model.User
		sql := "SELECT id, name, username, email, is_active from users where id = ?"
		if err := r.dbMasterRead.Raw(sql, id).First(&user).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		output <- model.Result{Data: user}
	}()
	return output
}

func (r *userManagementRepo) AddUserData(user model.NewUser) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		if err := r.dbMasterWrite.Create(&user).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		user.CreatedAtSTr = user.CreatedAt.Format(constant.DateTimeFormat)
		user.UpdatedAtStr = user.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: user}
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
		if err := r.dbMasterWrite.Model(&userTable).Where("id = ?", id).Update(updateData).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
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
		if err := r.dbMasterWrite.Where("id = ?", id).Delete(&user).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		output <- model.Result{}
	}()
	return output
}
