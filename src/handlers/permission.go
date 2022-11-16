package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func CreatePermissions(orgKey string, resourceId string, permissions *models.Permissions) error {
	if len(permissions.Permissions) < 1 {
		return errors.New("no permisiions")
	}
	for _, perm := range permissions.Permissions {
		perm.Actions = RemoveDuplicateActions(perm.Actions)
	}
	err := repositories.CreatePermissions(orgKey, resourceId, permissions)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDuplicateActions(values []models.Action) []models.Action {

	processed := make(map[models.Action]struct{})

	uniqValues := make([]models.Action, 0)
	for _, uid := range values {
		if _, ok := processed[uid]; ok {
			continue
		}
		uniqValues = append(uniqValues, uid)
		processed[uid] = struct{}{}
	}

	return uniqValues
}
