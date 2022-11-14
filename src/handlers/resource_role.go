package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResourceRoles(resource_id string, resource_roles *[]models.ResourceRole) error {
	return repositories.GetResourceRoles(resource_id, resource_roles)
}

func GetResourceRole(resource_id string, id string, resource_role *models.ResourceRole) error {
	return repositories.GetResourceRole(resource_id, id, resource_role)
}

func CreateResourceRole(resource_id string, resource_role *models.ResourceRole) error {
	return repositories.CreateResourceRole(resource_id, resource_role)
}

func DeleteResourceRole(resource_id string, id string) error {
	return repositories.DeleteResourceRole(resource_id, id)
}

func UpdateResourceRole(resource_id string, id string, resource_role *models.ResourceRole,
	reqResourceRole *models.ResourceRoleUpdateRequest) error {
	err := repositories.GetResourceRole(resource_id, id, resource_role)
	if err != nil {
		return err
	}
	resource_role.Name = reqResourceRole.Name
	return repositories.UpdateResourceRole(resource_role)
}

func CheckResourceRoleExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceRoleExistsById(id, &exists)
	return exists, err
}

func CheckResourceRoleExistsByKey(resource_id string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceRoleExistsByKey(resource_id, key, &exists)
	return exists, err
}
