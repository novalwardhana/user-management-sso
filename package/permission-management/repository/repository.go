package repository

import (
	"fmt"

	"github.com/novalwardhana/user-management-sso/config/postgres"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/permission-management/model"
)

type permissionManagementRepo struct {
	dbMaster *postgres.DBConnection
}

type PermissionManagementRepo interface {
	GetPermissionData(model.ListParams) <-chan model.Result
	GetPermissionByID(int) <-chan model.Result
	AddPermissionData(model.NewPermission) <-chan model.Result
	UpdatePermissionData(int, model.UpdatePermission) <-chan model.Result
	DeletePermissionData(int) <-chan model.Result
	GetTotalPermissionData(model.ListParams) <-chan model.Result
}

func NewPermissionManagementRepo(dbMaster *postgres.DBConnection) PermissionManagementRepo {
	return &permissionManagementRepo{
		dbMaster: dbMaster,
	}
}

func (r *permissionManagementRepo) GetPermissionData(params model.ListParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var permissions []model.Permission
		sql := `SELECT id, code, name, description FROM permissions order by id desc offset ? limit ?`
		rows, err := r.dbMaster.Read.Raw(sql, params.Offset, params.Limit).Rows()
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
		if err := r.dbMaster.Read.Raw(sql, id).First(&permission).Error; err != nil {
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
		if err := r.dbMaster.Read.Create(&permission).Error; err != nil {
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

		update := r.dbMaster.Write.Model(&permissionTable).Where("id = ?", id).Update(updateData)
		if update.Error != nil {
			output <- model.Result{Error: update.Error}
			return
		}
		if update.RowsAffected == 0 {
			output <- model.Result{Error: fmt.Errorf("Cannot update, permission data not found")}
			return
		}
		permission.ID = id
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

		delete := r.dbMaster.Write.Where("id = ?", id).Delete(&permission)
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

func (r *permissionManagementRepo) GetTotalPermissionData(params model.ListParams) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var count int
		sql := `select count(*) from permissions`
		if err := r.dbMaster.Read.Raw(sql).Count(&count).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}

		output <- model.Result{Data: count}
	}()
	return output
}
