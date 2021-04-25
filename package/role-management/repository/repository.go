package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/role-management/model"
	log "github.com/sirupsen/logrus"
)

type roleManagementRepo struct {
	dbMasterRead  *gorm.DB
	dbMasterWrite *gorm.DB
}

type RoleManagementRepo interface {
	GetRoleData() <-chan model.Result
	GetRoleByID(int) <-chan model.Result
	AddRoleData(model.NewRole) <-chan model.Result
	AddRolePermissionData([]model.RoleHasPermission) <-chan model.Result
	UpdateRoleData(int, model.UpdateRole) <-chan model.Result
	DeleteRoleData(int) <-chan model.Result
	DeleteRolePermissionData(int) <-chan model.Result
}

func NewRoleManagementRepo(dbMasterRead, dbMasterWrite *gorm.DB) RoleManagementRepo {
	return &roleManagementRepo{
		dbMasterRead:  dbMasterRead,
		dbMasterWrite: dbMasterWrite,
	}
}

func (r *roleManagementRepo) GetRoleData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var roles []model.Role
		sql := `SELECT id, code, name, "group", description from roles`
		rows, err := r.dbMasterRead.Raw(sql).Rows()
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
				&role.Description,
			); err != nil {
				output <- model.Result{Error: err}
				return
			} else {
				roles = append(roles, role)
			}
		}
		output <- model.Result{Data: roles}
	}()
	return output
}

func (r *roleManagementRepo) GetRoleByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var role model.Role
		sql := `SELECT id, code, name, "group", description from roles where id = ? `
		if err := r.dbMasterRead.Raw(sql, id).First(&role).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		output <- model.Result{Data: role}
	}()
	return output
}

func (r *roleManagementRepo) AddRoleData(user model.NewRole) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		create := r.dbMasterWrite.Create(&user)
		if create.Error != nil {
			output <- model.Result{Error: create.Error}
			return
		}

		user.CreatedAtStr = user.CreatedAt.Format(constant.DateTimeFormat)
		user.UpdatedAtStr = user.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: user}
	}()
	return output
}

func (r *roleManagementRepo) AddRolePermissionData(rolePermissions []model.RoleHasPermission) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		for _, rolePermission := range rolePermissions {
			create := r.dbMasterWrite.Create(&rolePermission)
			if create.Error != nil {
				log.Error(fmt.Sprintf("An error occured when insert role permission: %s\n", create.Error.Error()))
			}
		}

		output <- model.Result{}
	}()
	return output
}

func (r *roleManagementRepo) UpdateRoleData(id int, role model.UpdateRole) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var updateRole model.UpdateRole
		var updateData = map[string]interface{}{
			"code":        role.Code,
			"name":        role.Name,
			"group":       role.Group,
			"description": role.Description,
			"updated_at":  role.UpdatedAt,
		}

		update := r.dbMasterWrite.Model(&updateRole).Where("id = ?", id).Update(updateData)
		if update.Error != nil {
			output <- model.Result{Error: update.Error}
			return
		}
		if update.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot update, role data not found")}
			return
		}

		role.UpdatedAtSTr = role.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: role}
	}()
	return output
}

func (r *roleManagementRepo) DeleteRoleData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var role model.Role

		delete := r.dbMasterWrite.Where("id = ?", id).Delete(&role)
		if delete.Error != nil {
			output <- model.Result{Error: delete.Error}
			return
		}
		if delete.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot delete, role data not found")}
			return
		}

		output <- model.Result{}
	}()
	return output
}

func (r *roleManagementRepo) DeleteRolePermissionData(roleID int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var rolePermission model.RoleHasPermission
		delete := r.dbMasterWrite.Where("role_id = ?", roleID).Delete(&rolePermission)
		if delete.Error != nil {
			output <- model.Result{Error: delete.Error}
			return
		}

		output <- model.Result{}
	}()
	return output
}
