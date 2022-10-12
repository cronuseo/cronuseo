package repositories

import (
	"errors"
	"strconv"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resourceRoles *[]models.ResourceRole, res_id string) {
	config.DB.Model(&models.ResourceRole{}).Where("resource_id = ?", res_id).Find(&resourceRoles)
}

func GetResourceRole(resourceRole *models.ResourceRole, resrole_id string) {
	config.DB.Where("id = ?", resrole_id).First(&resourceRole)
}

func CreateResourceRoleAction(resourceRole *models.ResourceRole) {
	config.DB.Create(&resourceRole)
}

func DeleteResourceRole(resourceRole *models.ResourceRole, resrole_id string) {
	config.DB.Where("id = ?", resrole_id).Delete(&resourceRole)
}

func UpdateResourceRole(resourceRole *models.ResourceRole, reqResourceRole *models.ResourceRole, resrole_id string) {
	config.DB.Where("id = ?", resrole_id).First(&resourceRole)
	resourceRole.Name = reqResourceRole.Name
	resourceRole.Description = reqResourceRole.Description
	config.DB.Save(&resourceRole)
}

func DeleteAllResourceRoles(res_id string) {
	config.DB.Where("resource_id = ?", res_id).Delete(&models.ResourceRole{})
}

func AddUserToResourceRole(resrole_id string, user_id string) {
	roleuser := models.ResourceRoleToUser{}
	int_resrole_id, _ := strconv.Atoi(resrole_id)
	int_user_id, _ := strconv.Atoi(user_id)
	roleuser.ResourceRoleID = int_resrole_id
	roleuser.UserID = int_user_id
	config.DB.Create(roleuser)

}

func AddGroupToResourceRole(resrole_id string, group_id string) {
	rolegroup := models.ResourceRoleToGroup{}
	int_resrole_id, _ := strconv.Atoi(resrole_id)
	int_group_id, _ := strconv.Atoi(group_id)
	rolegroup.ResourceRoleID = int_resrole_id
	rolegroup.GroupID = int_group_id
	config.DB.Create(rolegroup)

}

func GetUResourceRoleWithGroupsAndUsers(resrole_id string, resourceRoleWithGroupsUsers *models.ResourceRoleWithGroupsUsers) {
	resourceRoleToGroup := []models.ResourceRoleToGroup{}
	resourceRoleToUser := []models.ResourceRoleToUser{}
	int_resrole_id, _ := strconv.Atoi(resrole_id)
	config.DB.Model(&models.ResourceRole{}).Select("id", "key", "name", "resource_id").Where("id = ?", resrole_id).Find(&resourceRoleWithGroupsUsers)
	config.DB.Model(&models.ResourceRoleToGroup{}).Where("resource_role_id = ?", int_resrole_id).Find(&resourceRoleToGroup)
	config.DB.Model(&models.ResourceRoleToUser{}).Where("resource_role_id = ?", int_resrole_id).Find(&resourceRoleToUser)
	if len(resourceRoleToUser) > 0 {
		for _, user := range resourceRoleToUser {
			user_id := user.UserID
			user := models.UserOnlyWithID{UserID: user_id}
			resourceRoleWithGroupsUsers.Users = append(resourceRoleWithGroupsUsers.Users, user)
		}
	}
	if len(resourceRoleToGroup) > 0 {
		for _, group := range resourceRoleToGroup {
			group_id := group.GroupID
			group := models.GroupOnlyWithID{GroupID: group_id}
			resourceRoleWithGroupsUsers.Groups = append(resourceRoleWithGroupsUsers.Groups, group)
		}
	}

}

func CheckResourceRoleExistsById(resact_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("id = ?", resact_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("resource role not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckResourceRoleExistsByKey(key string, res_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("key = ? AND resource_id = ?", key, res_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
