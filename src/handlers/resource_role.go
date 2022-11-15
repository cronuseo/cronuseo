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

func PatchResourceRole(resource_id string, resource_role_id string,
	resourceRole *models.ResourceRole, resourceRolePatch *models.ResourceRolePatchRequest) error {
	for _, operation := range resourceRolePatch.Operations {
		operation.Values = RemoveDuplicateValues(operation.Values)
	}
	err := repositories.PatchResourceRole(resource_role_id, resourceRolePatch)
	if err != nil {
		return err
	}
	return GetResourceRole(resource_id, resource_role_id, resourceRole)

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

func RemoveDuplicateValues(values []models.Value) []models.Value {

	processed := make(map[models.Value]struct{})

	uniqValues := make([]models.Value, 0)
	for _, uid := range values {
		if _, ok := processed[uid]; ok {
			continue
		}
		uniqValues = append(uniqValues, uid)
		processed[uid] = struct{}{}
	}

	return uniqValues
}
