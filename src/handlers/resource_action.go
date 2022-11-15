package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResourceActions(resource_id string, resourceActions *[]models.ResourceAction) error {
	return repositories.GetResourceActions(resource_id, resourceActions)
}

func GetResourceAction(resource_id string, id string, resourceAction *models.ResourceAction) error {
	return repositories.GetResourceAction(resource_id, id, resourceAction)
}

func CreateResourceAction(resource_id string, resourceAction *models.ResourceAction) error {
	return repositories.CreateResourceAction(resource_id, resourceAction)
}

func DeleteResourceAction(resource_id string, id string) error {
	return repositories.DeleteResourceAction(resource_id, id)
}

func UpdateResourceAction(resource_id string, id string, resourceAction *models.ResourceAction,
	reqResourceAction *models.ResourceActionUpdateRequest) error {
	err := repositories.GetResourceAction(resource_id, id, resourceAction)
	if err != nil {
		return err
	}
	resourceAction.Name = reqResourceAction.Name
	return repositories.UpdateResourceAction(resourceAction)
}

func CheckResourceActionExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceActionExistsById(id, &exists)
	return exists, err
}

func CheckResourceActionExistsByKey(resource_id string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceActionExistsByKey(resource_id, key, &exists)
	return exists, err
}
