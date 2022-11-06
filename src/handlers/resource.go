package handlers

import (
	"errors"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResources(resources *[]models.Resource, projId string) {
	repositories.GetResources(resources, projId)
}

func GetResource(resource *models.Resource, resId string) {
	repositories.GetResource(resource, resId)
}

func CreateResource(resource *models.Resource) {
	repositories.CreateResource(resource)
}

func DeleteResource(resource *models.Resource, resId string) {
	DeleteAllResourceActions(string(resId))
	repositories.DeleteAllResourceRoles(string(resId))
	repositories.DeleteResource(resource, resId)
}

func UpdateResource(resource *models.Resource, reqResource *models.Resource, resId string) {
	repositories.GetResource(resource, resId)
	resource.Name = reqResource.Name
	resource.Description = reqResource.Description
	repositories.UpdateResource(resource)
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
