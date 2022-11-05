package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resourceRoles *[]models.ResourceRole, resId string) {
	config.DB.Model(&models.ResourceRole{}).Where("resource_id = ?", resId).Find(&resourceRoles)
}

func GetResourceRole(resourceRole *models.ResourceRole, resroleId string) {
	config.DB.Where("id = ?", resroleId).First(&resourceRole)
}

func CreateResourceRole(resourceRole *models.ResourceRole) {
	config.DB.Create(&resourceRole)
}

func DeleteResourceRole(resourceRole *models.ResourceRole, resroleId string) {
	config.DB.Where("id = ?", resroleId).Delete(&resourceRole)
}

func UpdateResourceRole(resourceRole *models.ResourceRole) {
	config.DB.Save(&resourceRole)
}

func DeleteAllResourceRoles(resId string) {
	config.DB.Where("resource_id = ?", resId).Delete(&models.ResourceRole{})
}

func AddUserToResourceRole(roleuser *models.ResourceRoleToUser) {
	config.DB.Create(roleuser)

}

func AddGroupToResourceRole(rolegroup *models.ResourceRoleToGroup) {
	config.DB.Create(rolegroup)

}

func AddResourceActionToResourceRole(roleaction *models.ResourceRoleToResourceAction, roleActionKey *models.ResourceRoleToResourceActionKey) {
	config.DB.Create(roleaction)
	config.DB.Create(roleActionKey)

}

func CheckResourceActionAlreadyAdded(resId string, resRoleId string, resActId string, exists bool) error {
	return config.DB.Model(&models.ResourceRoleToResourceAction{}).Select("count(*) > 0").Where("resource_id = ? AND resource_role_id = ? AND resource_action_id = ?", resId, resRoleId, resActId).Find(&exists).Error
}

func CheckGroupAlreadyAdded(resRoleId string, groupId string, exists bool) error {
	return config.DB.Model(&models.ResourceRoleToGroup{}).Where("resource_role_id = ? AND group_id = ?", resRoleId, groupId).Find(&exists).Error

}

func CheckUserAlreadyAdded(resRoleId string, userId string, exists bool) error {
	return config.DB.Model(&models.ResourceRoleToUser{}).Where("resource_role_id = ? AND user_id = ?", resRoleId, userId).Find(&exists).Error
}

func GetUResourceRoleWithGroupsAndUsers(resroleId string, resourceRoleWithGroupsUsers *models.ResourceRoleWithGroupsUsers, resourceRoleToGroup *[]models.ResourceRoleToGroup, resourceRoleToUser *[]models.ResourceRoleToUser, resourceRoleToAction *[]models.ResourceRoleToResourceAction) {
	config.DB.Model(&models.ResourceRole{}).Select("id", "key", "name", "resource_id").Where("id = ?", resroleId).Find(&resourceRoleWithGroupsUsers)
	config.DB.Model(&models.ResourceRoleToGroup{}).Where("resource_role_id = ?", resroleId).Find(&resourceRoleToGroup)
	config.DB.Model(&models.ResourceRoleToUser{}).Where("resource_role_id = ?", resroleId).Find(&resourceRoleToUser)
	config.DB.Model(&models.ResourceRoleToResourceAction{}).Where("resource_role_id = ?", resroleId).Find(&resourceRoleToAction)
}

func CheckResourceRoleExistsById(resRoleId string, exists bool) error {
	return config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("id = ?", resRoleId).Find(&exists).Error
}

func CheckResourceRoleExistsByKey(key string, resId string, exists bool) error {
	return config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("key = ? AND resource_id = ?", key, resId).Find(&exists).Error
}
