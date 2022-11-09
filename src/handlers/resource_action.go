package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/repositories"

	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceActions(resourceActions *[]models.ResourceAction, resId string) error {
	return repositories.GetResourceActions(resourceActions, resId)
}

func GetResourceAction(resourceAction *models.ResourceAction, resactId string) error {
	return repositories.GetResourceAction(resourceAction, resactId)
}

func CreateResourceAction(resourceAction *models.ResourceAction) error {
	return repositories.CreateResourceAction(resourceAction)
}

func DeleteResourceAction(resourceAction *models.ResourceAction, resactId string) error {
	return repositories.DeleteResourceAction(resourceAction, resactId)
}

func UpdateResourceAction(resourceAction *models.ResourceAction, reqResourceAction *models.ResourceAction,
	resactId string) error {
	err := repositories.GetResourceAction(resourceAction, resactId)
	if err != nil {
		return err
	}
	resourceAction.Name = reqResourceAction.Name
	resourceAction.Description = reqResourceAction.Description
	return repositories.UpdateResourceAction(resourceAction)
}

func CheckResourceActionExistsById(resactId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceActionExistsById(resactId, &exists)
	if err != nil {
		return false, errors.New("resource action not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckResourceActionExistsByKey(key string, resId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceActionExistsByKey(key, resId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
