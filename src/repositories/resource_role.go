package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resourceRoles *[]models.ResourceRole, resId string) error {
	return config.DB.Model(&models.ResourceRole{}).Where("resource_id = ?", resId).Find(&resourceRoles).Error
}

func GetResourceRole(resourceRole *models.ResourceRole, resroleId string) error {
	return config.DB.Where("id = ?", resroleId).First(&resourceRole).Error
}

func CreateResourceRole(resourceRole *models.ResourceRole) error {
	return config.DB.Create(&resourceRole).Error
}

func DeleteResourceRole(resourceRole *models.ResourceRole, resroleId string) error {
	return config.DB.Where("id = ?", resroleId).Delete(&resourceRole).Error
}

func UpdateResourceRole(resourceRole *models.ResourceRole) error {
	return config.DB.Save(&resourceRole).Error
}

func DeleteAllResourceRoles(resId string) error {
	return config.DB.Where("resource_id = ?", resId).Delete(&models.ResourceRole{}).Error
}

func AddUserToResourceRole(roleuser *models.ResourceRoleToUser) error {
	return config.DB.Create(roleuser).Error

}

func AddGroupToResourceRole(rolegroup *models.ResourceRoleToGroup) error {
	return config.DB.Create(rolegroup).Error

}

func AddResourceActionToResourceRole(roleaction *models.ResourceRoleToResourceAction,
	roleActionKey *models.ResourceRoleToResourceActionKey) error {
	err := config.DB.Create(roleaction).Error
	if err != nil {
		return err
	}
	return config.DB.Create(roleActionKey).Error

}

func CheckResourceActionAlreadyAdded(resId string, resRoleId string, resActId string, exists *bool) error {
	return config.DB.Model(&models.ResourceRoleToResourceAction{}).Select("count(*) > 0").Where(
		"resource_id = ? AND resource_role_id = ? AND resource_action_id = ?", resId, resRoleId,
		resActId).Find(exists).Error
}

func CheckGroupAlreadyAdded(resRoleId string, groupId string, exists *bool) error {
	return config.DB.Model(&models.ResourceRoleToGroup{}).Select("count(*) > 0").Where(
		"resource_role_id = ? AND group_id = ?", resRoleId, groupId).Find(exists).Error

}

func CheckUserAlreadyAdded(resRoleId string, userId string, exists *bool) error {
	return config.DB.Model(&models.ResourceRoleToUser{}).Select("count(*) > 0").Where(
		"resource_role_id = ? AND user_id = ?", resRoleId, userId).Find(exists).Error
}

func GetUResourceRoleWithGroupsAndUsers(resRoleId string,
	resourceRoleWithGroupsUsers *models.ResourceRoleWithGroupsUsers, resourceRoleToGroup *[]models.ResourceRoleToGroup,
	resourceRoleToUser *[]models.ResourceRoleToUser, resourceRoleToAction *[]models.ResourceRoleToResourceAction) error {
	err := config.DB.Model(&models.ResourceRole{}).Select("id", "key", "name", "resource_id").Where(
		"id = ?", resRoleId).Find(&resourceRoleWithGroupsUsers).Error
	if err != nil {
		return err
	}
	err = config.DB.Model(&models.ResourceRoleToGroup{}).Where("resource_role_id = ?",
		resRoleId).Find(&resourceRoleToGroup).Error
	if err != nil {
		return err
	}
	err = config.DB.Model(&models.ResourceRoleToUser{}).Where("resource_role_id = ?",
		resRoleId).Find(&resourceRoleToUser).Error
	if err != nil {
		return err
	}
	return config.DB.Model(&models.ResourceRoleToResourceAction{}).Where("resource_role_id = ?",
		resRoleId).Find(&resourceRoleToAction).Error
}

func CheckResourceRoleExistsById(resRoleId string, exists *bool) error {
	return config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("id = ?", resRoleId).Find(exists).Error
}

func CheckResourceRoleExistsByKey(key string, resId string, exists *bool) error {
	return config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where(
		"key = ? AND resource_id = ?", key, resId).Find(exists).Error
}
