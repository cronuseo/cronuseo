package handlers

import (
	"errors"
	"strconv"

	"github.com/shashimalcse/Cronuseo/repositories"

	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resourceRoles *[]models.ResourceRole, resId string) error {
	return repositories.GetResourceRoles(resourceRoles, resId)
}

func GetResourceRole(resourceRole *models.ResourceRole, resRoleId string) error {
	return repositories.GetResourceRole(resourceRole, resRoleId)
}

func CreateResourceRoleAction(resourceRole *models.ResourceRole) error {
	return repositories.CreateResourceRole(resourceRole)
}

func DeleteResourceRole(resourceRole *models.ResourceRole, resRoleId string) error {
	return repositories.DeleteResourceRole(resourceRole, resRoleId)
}

func UpdateResourceRole(resourceRole *models.ResourceRole, reqResourceRole *models.ResourceRole,
	resroleId string) error {
	err := repositories.GetResourceRole(resourceRole, resroleId)
	if err != nil {
		return err
	}
	resourceRole.Name = reqResourceRole.Name
	resourceRole.Description = reqResourceRole.Description
	return repositories.UpdateResourceRole(resourceRole)
}

func DeleteAllResourceRoles(resId string) error {
	return repositories.DeleteAllResourceRoles(resId)
}

func AddUserToResourceRole(resRoleId string, userId string) error {
	roleUser := models.ResourceRoleToUser{}
	intResRoleId, _ := strconv.Atoi(resRoleId)
	intUserId, _ := strconv.Atoi(userId)
	roleUser.ResourceRoleID = intResRoleId
	roleUser.UserID = intUserId
	return repositories.AddUserToResourceRole(&roleUser)

}

func AddGroupToResourceRole(resRoleId string, groupId string) error {
	roleGroup := models.ResourceRoleToGroup{}
	intResRoleId, _ := strconv.Atoi(resRoleId)
	intGroupId, _ := strconv.Atoi(groupId)
	roleGroup.ResourceRoleID = intResRoleId
	roleGroup.GroupID = intGroupId
	return repositories.AddGroupToResourceRole(&roleGroup)

}

func AddResourceActionToResourceRole(resId string, resRoleId string, resActId string) error {
	roleAction := models.ResourceRoleToResourceAction{}
	intResRoleId, _ := strconv.Atoi(resRoleId)
	intResActId, _ := strconv.Atoi(resActId)
	intResId, _ := strconv.Atoi(resId)
	roleAction.ResourceRoleID = intResRoleId
	roleAction.ResourceActionID = intResActId
	roleAction.ResourceID = intResId
	var resourceRole *models.ResourceRole
	err := repositories.GetResourceRole(resourceRole, resRoleId)
	if err != nil {
		return err
	}
	var resourceAction *models.ResourceAction
	err = repositories.GetResourceAction(resourceAction, resActId)
	if err != nil {
		return err
	}
	var resource *models.Resource
	err = repositories.GetResource(resource, resId)
	if err != nil {
		return err
	}
	resKey := resource.Key
	resRoleKey := resourceRole.Key
	resActKey := resourceAction.Key
	roleActionKey := models.ResourceRoleToResourceActionKey{Resource: resKey, ResourceAction: resActKey,
		ResourceRole: resRoleKey}
	return repositories.AddResourceActionToResourceRole(&roleAction, &roleActionKey)

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
	resourceRoleWithGroupsUsers *models.ResourceRoleWithGroupsUsers) error {
	resourceRoleToGroup := []models.ResourceRoleToGroup{}
	resourceRoleToUser := []models.ResourceRoleToUser{}
	resourceRoleToAction := []models.ResourceRoleToResourceAction{}
	err := repositories.GetUResourceRoleWithGroupsAndUsers(resrole_id, resourceRoleWithGroupsUsers,
		&resourceRoleToGroup, &resourceRoleToUser, &resourceRoleToAction)
	if err != nil {
		return err
	}
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
	return nil

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
