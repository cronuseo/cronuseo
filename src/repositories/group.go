package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetGroups(groups *[]models.Group, orgId string) {
	config.DB.Model(&models.Group{}).Where("organization_id = ?", orgId).Find(&groups)
}

func GetGroup(group *models.Group, groupId string) {
	config.DB.Where("id = ?", groupId).First(&group)
}

func CreateGroup(group *models.Group) {
	config.DB.Create(&group)
}

func DeleteGroup(group *models.Group, groupId string) {
	config.DB.Where("id = ?", groupId).Delete(&group)
}

func UpdateGroup(group *models.Group) {
	config.DB.Save(&group)
}

func AddUserToGroup(groupuser models.GroupUser) {
	config.DB.Create(groupuser)

}

func GetUsersFromGroup(groupId int, resGroupUsers *models.GroupUsers, groupusers []models.GroupUser) {
	config.DB.Model(&models.Group{}).Select("id", "key", "name", "organization_id").Where("id = ?", groupId).Find(&resGroupUsers)
	config.DB.Model(&models.GroupUser{}).Where("group_id = ?", groupId).Find(&groupusers)

}

func CheckGroupExistsById(groupId string, exists bool) error {
	return config.DB.Model(&models.Group{}).Select("count(*) > 0").Where("id = ?", groupId).Find(&exists).Error
}

func CheckGroupExistsByKey(key string, org_id string, exists bool) error {
	return config.DB.Model(&models.Group{}).Select("count(*) > 0").Where("key = ? AND organization_id = ?", key, org_id).Find(&exists).Error
}
