package repositories

import (
	"errors"
	"strconv"

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

func AddUserToGroup(group_id string, user_id string) {
	groupuser := models.GroupUser{}
	int_group_id, _ := strconv.Atoi(group_id)
	int_user_id, _ := strconv.Atoi(user_id)
	groupuser.GroupID = int_group_id
	groupuser.UserID = int_user_id
	config.DB.Create(groupuser)

}

func GetUsersFromGroup(group_id string, resGroupusers *models.GroupUsers) {
	groupusers := []models.GroupUser{}
	int_group_id, _ := strconv.Atoi(group_id)
	config.DB.Model(&models.Group{}).Select("id", "key", "name", "organization_id").Where("id = ?", group_id).Find(&resGroupusers)
	config.DB.Model(&models.GroupUser{}).Where("group_id = ?", int_group_id).Find(&groupusers)
	if len(groupusers) > 0 {
		for _, groupuser := range groupusers {
			user_id := groupuser.UserID
			user := models.UserOnlyWithID{UserID: user_id}
			resGroupusers.Users = append(resGroupusers.Users, user)
		}
	}

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
