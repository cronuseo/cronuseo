package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResources(resources *[]models.Resource, projId string) error {
	return repositories.GetResources(resources, projId)
}

func GetResource(resource *models.Resource, resId string) error {
	return repositories.GetResource(resource, resId)
}

func CreateResource(resource *models.Resource) error {
	return repositories.CreateResource(resource)
}

func DeleteResource(resource *models.Resource, resId string) error {
	err := repositories.DeleteAllResourceActions(string(resId))
	if err != nil {
		return err
	}
	err = repositories.DeleteAllResourceRoles(string(resId))
	if err != nil {
		return err
	}
	return repositories.DeleteResource(resource, resId)
}

func UpdateResource(resource *models.Resource, reqResource *models.Resource, resId string) error {
	err := repositories.GetResource(resource, resId)
	if err != nil {
		return err
	}
	resource.Name = reqResource.Name
	resource.Description = reqResource.Description
	return repositories.UpdateResource(resource)
}

func CheckResourceExistsById(resId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceExistsById(resId, &exists)
	if err != nil {
		return false, errors.New("resource not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckResourceExistsByKey(key string, projId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceExistsByKey(key, projId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
