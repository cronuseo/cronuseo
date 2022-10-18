package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func Check(resource string, role string, action string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRoleToResourceActionKey{}).Select("count(*) > 0").Where("resource = ? AND resource_role = ? AND resource_action = ?", resource, role, action).Find(&exists).Error
	if err != nil {
		return false, errors.New("not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
