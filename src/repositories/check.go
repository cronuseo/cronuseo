package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func Check(resource string, role string, action string, exists *bool) error {
	return config.DB.Model(&models.ResourceRoleToResourceActionKey{}).Select(
		"count(*) > 0").Where("resource = ? AND resource_role = ? AND resource_action = ?", resource, role,
		action).Find(exists).Error
}
