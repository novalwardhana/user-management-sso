package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/novalwardhana/user-management-sso/package/user-management/model"
)

type userManagementRepo struct {
	dbMasterRead  *gorm.DB
	dbMasterWrite *gorm.DB
}

type UserManagementRepo interface {
	GetUserData() <-chan model.Result
	GetUserByID(int) <-chan model.Result
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
		sql := "SELECT id, username, email FROM users"
		rows, err := r.dbMasterRead.Raw(sql).Rows()
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		for rows.Next() {
			user := model.User{}
			err := rows.Scan(
				&user.ID,
				&user.Username,
				&user.Email,
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
		sql := "SELECT id, username, email from users where id = ?"
		if err := r.dbMasterRead.Raw(sql, id).First(&user).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		output <- model.Result{Data: user}
	}()
	return output
}
