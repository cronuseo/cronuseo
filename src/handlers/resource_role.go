package handlers

import (
	"errors"
	"github.com/shashimalcse/Cronuseo/repositories"
	"strconv"

	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resourceRoles *[]models.ResourceRole, resId string) {
	repositories.GetResourceRoles(resourceRoles, resId)
}

func GetResourceRole(resourceRole *models.ResourceRole, resRoleId string) {
	repositories.GetResourceRole(resourceRole, resRoleId)
}

func CreateResourceRoleAction(resourceRole *models.ResourceRole) {
	repositories.CreateResourceRole(resourceRole)
}

func DeleteResourceRole(resourceRole *models.ResourceRole, resRoleId string) {
	repositories.DeleteResourceRole(resourceRole, resRoleId)
}

func UpdateResourceRole(resourceRole *models.ResourceRole, reqResourceRole *models.ResourceRole, resroleId string) {
	repositories.GetResourceRole(resourceRole, resroleId)
	resourceRole.Name = reqResourceRole.Name
	resourceRole.Description = reqResourceRole.Description
	repositories.UpdateResourceRole(resourceRole)
}

func DeleteAllResourceRoles(resId string) {
	repositories.DeleteAllResourceRoles(resId)
}

func AddUserToResourceRole(resRoleId string, userId string) {
	roleUser := models.ResourceRoleToUser{}
	intResRoleId, _ := strconv.Atoi(resRoleId)
	intUserId, _ := strconv.Atoi(userId)
	roleUser.ResourceRoleID = intResRoleId
	roleUser.UserID = intUserId
	repositories.AddUserToResourceRole(&roleUser)

}

func AddGroupToResourceRole(resRoleId string, groupId string) {
	roleGroup := models.ResourceRoleToGroup{}
	intResRoleId, _ := strconv.Atoi(resRoleId)
	intGroupId, _ := strconv.Atoi(groupId)
	roleGroup.ResourceRoleID = intResRoleId
	roleGroup.GroupID = intGroupId
	repositories.AddGroupToResourceRole(&roleGroup)

}

func AddResourceActionToResourceRole(resId string, resRoleId string, resActId string) {
	roleAction := models.ResourceRoleToResourceAction{}
	intResRoleId, _ := strconv.Atoi(resRoleId)
	intResActId, _ := strconv.Atoi(resActId)
	intResId, _ := strconv.Atoi(resId)
	roleAction.ResourceRoleID = intResRoleId
	roleAction.ResourceActionID = intResActId
	roleAction.ResourceID = intResId
	var resourceRole *models.ResourceRole
	repositories.GetResourceRole(resourceRole, resRoleId)
	var resourceAction *models.ResourceAction
	repositories.GetResourceAction(resourceAction, resActId)
	var resource *models.Resource
	repositories.GetResource(resource, resId)
	resKey := resource.Key
	resRoleKey := resourceRole.Key
	resActKey := resourceAction.Key
	roleActionKey := models.ResourceRoleToResourceActionKey{Resource: resKey, ResourceAction: resActKey,
		ResourceRole: resRoleKey}
	repositories.AddResourceActionToResourceRole(&roleAction, &roleActionKey)

}

func CheckResourceActionAlreadyAdded(resId string, resRoleId string, resActId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceActionAlreadyAdded(resId, resRoleId, resActId, &exists)
	if err != nil {
		return false, errors.New("not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckGroupAlreadyAdded(resRoleId string, groupId string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupAlreadyAdded(resRoleId, groupId, &exists)
	if err != nil {
		return false, errors.New("not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckUserAlreadyAdded(resRoleId string, userId string) (bool, error) {
	var exists bool
	err := repositories.CheckUserAlreadyAdded(resRoleId, userId, &exists)
	if err != nil {
		return false, errors.New("not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func GetUResourceRoleWithGroupsAndUsers(resrole_id string,
	resourceRoleWithGroupsUsers *models.ResourceRoleWithGroupsUsers) {
	resourceRoleToGroup := []models.ResourceRoleToGroup{}
	resourceRoleToUser := []models.ResourceRoleToUser{}
	resourceRoleToAction := []models.ResourceRoleToResourceAction{}
	repositories.GetUResourceRoleWithGroupsAndUsers(resrole_id, resourceRoleWithGroupsUsers,
		&resourceRoleToGroup, &resourceRoleToUser, &resourceRoleToAction)
	if len(resourceRoleToUser) > 0 {
		for _, user := range resourceRoleToUser {
			userId := user.UserID
			user := models.UserOnlyWithID{UserID: userId}
			resourceRoleWithGroupsUsers.Users = append(resourceRoleWithGroupsUsers.Users, user)
		}
	}
	if len(resourceRoleToGroup) > 0 {
		for _, group := range resourceRoleToGroup {
			groupId := group.GroupID
			group := models.GroupOnlyWithID{GroupID: groupId}
			resourceRoleWithGroupsUsers.Groups = append(resourceRoleWithGroupsUsers.Groups, group)
		}
	}
	if len(resourceRoleToAction) > 0 {
		for _, action := range resourceRoleToAction {
			actionId := action.ResourceActionID
			action := models.ResourceActionWithID{ResourceActionID: actionId}
			resourceRoleWithGroupsUsers.ResourceActions = append(resourceRoleWithGroupsUsers.ResourceActions, action)
		}
	}

}

func CheckResourceRoleExistsById(resRoleId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceRoleExistsById(resRoleId, &exists)
	if err != nil {
		return false, errors.New("resource role not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckResourceRoleExistsByKey(key string, resId string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceRoleExistsByKey(key, resId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
