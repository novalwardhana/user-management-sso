package repository

import (
	"encoding/json"
	"fmt"

	"github.com/novalwardhana/user-management-sso/config/postgres"
	"github.com/novalwardhana/user-management-sso/global/constant"
	"github.com/novalwardhana/user-management-sso/package/role-management/model"
	log "github.com/sirupsen/logrus"
)

type roleManagementRepo struct {
	dbMaster *postgres.DBConnection
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

func NewRoleManagementRepo(dbMaster *postgres.DBConnection) RoleManagementRepo {
	return &roleManagementRepo{
		dbMaster: dbMaster,
	}
}

func (r *roleManagementRepo) GetRoleData() <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)
		var listRoles []model.ListRole

		sql := `SELECT 
				r.id, 
				r.code, 
				r.name, 
				r."group", 
				r.description,
				concat('[', string_agg('{"code":"' || p.code ||'","name":"' || p.name || '"}', ','), ']') as permissions 
			from roles r
			left join role_has_permissions rhp on r.id = rhp.role_id
			left join permissions p on rhp.permission_id = p.id
			group by r.id order by r.id desc`
		rows, err := r.dbMaster.Read.Raw(sql).Rows()
		if err != nil {
			output <- model.Result{Error: err}
			return
		}
		for rows.Next() {
			listRole := model.ListRole{}
			if err := rows.Scan(
				&listRole.ID,
				&listRole.Code,
				&listRole.Name,
				&listRole.Group,
				&listRole.Description,
				&listRole.Permissions,
			); err != nil {
				output <- model.Result{Error: err}
				return
			} else {
				permissions := []map[string]interface{}{}
				_ = json.Unmarshal([]byte(listRole.Permissions), &permissions)
				listRole.PermissionArrays = permissions
				listRoles = append(listRoles, listRole)
			}
		}
		output <- model.Result{Data: listRoles}
	}()
	return output
}

func (r *roleManagementRepo) GetRoleByID(id int) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		var listRole model.ListRole
		sql := `SELECT 
				r.id, 
				r.code, 
				r.name, 
				r."group",  
				r.description,
				concat('[', string_agg('{"id":' || p.id ||','|| '"code":"' || p.code ||'","name":"' || p.name || '"}', ','), ']') as permissions 
			from roles r
			left join role_has_permissions rhp on r.id = rhp.role_id
			left join permissions p on rhp.permission_id = p.id
			where r.id = ?
			group by r.id order by r.id desc limit 1`
		if err := r.dbMaster.Read.Raw(sql, id).Scan(&listRole).Error; err != nil {
			output <- model.Result{Error: err}
			return
		}

		permissions := []map[string]interface{}{}
		if err := json.Unmarshal([]byte(listRole.Permissions), &permissions); err != nil {
			output <- model.Result{Error: err}
			return
		}
		listRole.PermissionArrays = permissions
		output <- model.Result{Data: listRole}
	}()
	return output
}

func (r *roleManagementRepo) AddRoleData(user model.NewRole) <-chan model.Result {
	output := make(chan model.Result)
	go func() {
		defer close(output)

		create := r.dbMaster.Write.Create(&user)
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
			create := r.dbMaster.Write.Create(&rolePermission)
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

		update := r.dbMaster.Write.Model(&updateRole).Where("id = ?", id).Update(updateData)
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

		delete := r.dbMaster.Write.Where("id = ?", id).Delete(&role)
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
		delete := r.dbMaster.Write.Where("role_id = ?", roleID).Delete(&rolePermission)
		if delete.Error != nil {
			output <- model.Result{Error: delete.Error}
			return
		}

		output <- model.Result{}
	}()
	return output
}
