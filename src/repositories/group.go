package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetGroups(groups *[]models.Group, org_id string) {
	config.DB.Model(&models.Group{}).Where("organization_id = ?", org_id).Find(&groups)
}

func GetGroup(group *models.Group, group_id string) {
	config.DB.Where("id = ?", group_id).First(&group)
}

func CreateGroup(group *models.Group) {
	config.DB.Create(&group)
}

func DeleteGroup(group *models.Group, group_id string) {
	DeleteAllResources(group_id)
	config.DB.Where("id = ?", group_id).Delete(&group)
}

func UpdateGroup(group *models.Group, reqGroup *models.Group, group_id string) {
	config.DB.Where("id = ?", group_id).First(&group)
	group.Name = reqGroup.Name
	group.Key = reqGroup.Key
	config.DB.Save(&group)
}

func CheckGroupExistsById(group_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Group{}).Select("count(*) > 0").Where("id = ?", group_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("group not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckGroupExistsByKey(key string, org_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Group{}).Select("count(*) > 0").Where("key = ? AND organization_id = ?", key, org_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
