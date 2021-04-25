package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/permission-management/model"
)

type permissionManagementRepo struct {
	dbMasterRead  *gorm.DB
	dbMasterWrite *gorm.DB
}

type PermissionManagementRepo interface {
	GetPermissionData() <-chan model.Result
	GetPermissionByID(int) <-chan model.Result
	AddPermissionData(model.NewPermission) <-chan model.Result
	UpdatePermissionData(int, model.UpdatePermission) <-chan model.Result
	DeletePermissionData(int) <-chan model.Result
}

func NewPermissionManagementRepo(dbMasterRead, dbMasterWrite *gorm.DB) PermissionManagementRepo {
	return &permissionManagementRepo{
		dbMasterRead:  dbMasterRead,
		dbMasterWrite: dbMasterWrite,
	}
}

func (r *permissionManagementRepo) GetPermissionData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var permissions []model.Permission
		sql := `SELECT id, code, name, description FROM permissions`
		rows, err := r.dbMasterRead.Raw(sql).Rows()
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
				&permission.Description,
			); err != nil {
				output <- model.Result{Error: err}
				return
			} else {
				permissions = append(permissions, permission)
			}
		}
		output <- model.Result{Data: permissions}
	}()
	return output
}

func (r *permissionManagementRepo) GetPermissionByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var permission model.Permission
		sql := `SELECT id, code, name, description FROM permissions WHERE id = ?`
		if err := r.dbMasterRead.Raw(sql, id).First(&permission).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		output <- model.Result{Data: permission}
	}()
	return output
}

func (r *permissionManagementRepo) AddPermissionData(permission model.NewPermission) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		if err := r.dbMasterRead.Create(&permission).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}
		permission.CreatedAtStr = permission.CreatedAt.Format(constant.DateTimeFormat)
		permission.UpdatedAtStr = permission.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: permission}
	}()
	return output
}

func (r *permissionManagementRepo) UpdatePermissionData(id int, permission model.UpdatePermission) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var permissionTable model.UpdatePermission
		var updateData = map[string]interface{}{
			"code":        permission.Code,
			"name":        permission.Name,
			"description": permission.Description,
			"updated_at":  permission.UpdatedAt,
		}

		update := r.dbMasterWrite.Model(&permissionTable).Where("id = ?", id).Update(updateData)
		if update.Error != nil {
			output <- model.Result{Error: update.Error}
			return
		}
		if update.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot update, permission data not found")}
			return
		}
		permission.UpdatedAtStr = permission.UpdatedAt.Format(constant.DateTimeFormat)
		output <- model.Result{Data: permission}
	}()
	return output
}

func (r *permissionManagementRepo) DeletePermissionData(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var permission model.Permission

		delete := r.dbMasterWrite.Where("id = ?", id).Delete(&permission)
		if delete.Error != nil {
			output <- model.Result{Error: delete.Error}
			return
		}
		if delete.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot delete, permission data not found")}
			return
		}
		output <- model.Result{}
	}()
	return output
}
