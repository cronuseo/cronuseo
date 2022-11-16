package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func CreatePermissions(orgKey string, resourceId string, permissions *models.Permissions) error {

	tx, err := config.DB.Begin()

	if err != nil {
		return err
	}
	// add resource role
	if len(permissions.Permissions) > 0 {
		stmt, err := tx.Prepare(
			`INSERT INTO resource_action_role(resource_id,resource_role_id,resource_action_id) VALUES($1, $2, $3)`)

		if err != nil {
			return err
		}

		defer stmt.Close()
		for _, permission := range permissions.Permissions {
			roleId := permission.RoleID
			var exists bool
			role_err := CheckResourceRoleExistsById(roleId, &exists)
			if !exists || role_err != nil {
				continue
			}
			for _, action := range permission.Actions {
				var exists bool
				action_err := CheckResourceActionExistsById(action.ActionID, &exists)
				if !exists || action_err != nil {
					continue
				}
				er := ISPermissionAlreadyExists(resourceId, roleId, action.ActionID, &exists)
				if exists || er != nil {
					continue
				}
				_, err = stmt.Exec(resourceId, roleId, action.ActionID)
				if err != nil {
					return err
				}
				var resourceKey string
				key_err := GetResourceKeyById(resourceId, &resourceKey)
				if key_err != nil {
					return key_err
				}
				var projectId string
				key_err = GetProjectIDById(resourceId, &projectId)
				if key_err != nil {
					return key_err
				}
				var projectKey string
				key_err = GetProjectKeyById(projectId, &projectKey)
				if key_err != nil {
					return key_err
				}
				var roleKey string
				key_err = GetRoleKeyById(roleId, &roleKey)
				if key_err != nil {
					return key_err
				}
				var actionKey string
				key_err = GetActionKeyById(action.ActionID, &actionKey)
				if key_err != nil {
					return key_err
				}
				value_key := orgKey + ":" + projectKey + ":" + resourceKey + ":" + roleKey + ":" + actionKey

				stmt2, err := tx.Prepare(
					`INSERT INTO resource_action_role_key(value_key) VALUES($1)`)

				if err != nil {
					return err
				}
				_, err = stmt2.Exec(value_key)
				if err != nil {
					return err
				}
			}
		}

	}

	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}
	return nil
}

func ISPermissionAlreadyExists(resourceId string, roleId string, actionID string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT resource_id FROM resource_action_role"+
			" WHERE resource_id = $1 AND resource_role_id = $2 AND resource_action_id = $3)",
		resourceId, roleId, actionID).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
