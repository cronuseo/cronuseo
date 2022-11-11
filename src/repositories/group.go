package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetGroups(groups *[]models.Group, orgId string) error {
	return config.DB.Model(&models.Group{}).Where("organization_id = ?", orgId).Find(&groups).Error
}

func GetGroup(group *models.Group, groupId string) error {
	return config.DB.Where("id = ?", groupId).First(&group).Error
}

func CreateGroup(group *models.Group) error {
	return config.DB.Create(&group).Error
}

func DeleteGroup(group *models.Group, groupId string) error {
	return config.DB.Where("id = ?", groupId).Delete(&group).Error
}

func UpdateGroup(group *models.Group) error {
	return config.DB.Save(&group).Error
}

func AddUserToGroup(groupuser models.GroupUser) error {

	return config.DB.Create(groupuser).Error

}

func GetUsersFromGroup(groupId int, resGroupUsers *models.GroupUsers, groupusers *[]models.GroupUser) error {
	err := config.DB.Model(&models.Group{}).Select("id", "key", "name", "organization_id").Where("id = ?",
		groupId).Find(&resGroupUsers).Error
	if err != nil {
		return err
	}
	return config.DB.Model(&models.GroupUser{}).Where("group_id = ?", groupId).Find(&groupusers).Error

}

func CheckGroupExistsById(groupId string, exists *bool) error {
	return config.DB.Model(&models.Group{}).Select("count(*) > 0").Where("id = ?", groupId).Find(exists).Error
}

func CheckGroupExistsByKey(key string, orgId string, exists *bool) error {
	return config.DB.Model(&models.Group{}).Select("count(*) > 0").Where("key = ? AND organization_id = ?",
		key, orgId).Find(exists).Error
}

func CheckGroupAlreadyInGroup(groupId string, userId string, exists *bool) error {
	return config.DB.Model(&models.GroupUser{}).Select(
		"count(*) > 0").Where("group_id = ? AND user_id = ?", groupId, userId).Find(exists).Error
}
