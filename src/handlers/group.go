package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetGroups(groups *[]models.Group, orgId string) error {
	return repositories.GetGroups(groups, orgId)
}

func GetGroup(group *models.Group, groupId string) error {
	return repositories.GetGroup(group, groupId)
}

func CreateGroup(group *models.Group) error {
	return repositories.CreateGroup(group)
}

func DeleteGroup(group *models.Group, groupId string) error {
	//todo delete group from role group
	return repositories.DeleteGroup(group, groupId)
}

func UpdateGroup(group *models.Group, reqGroup *models.Group, groupId string) error {
	err := repositories.GetGroup(group, groupId)
	if err != nil {
		return err
	}
	group.Name = reqGroup.Name
	group.Key = reqGroup.Key
	return repositories.UpdateGroup(group)
}

func AddUserToGroup(groupId string, userId string) error {
	groupuser := models.GroupUser{}
	intGroupId, _ := strconv.Atoi(groupId)
	intUserId, _ := strconv.Atoi(userId)
	groupuser.GroupID = intGroupId
	groupuser.UserID = intUserId
	return repositories.AddUserToGroup(groupuser)

}

func GetUsersFromGroup(groupId string, resGroupusers *models.GroupUsers) error {
	var groupusers []models.GroupUser
	intGroupId, _ := strconv.Atoi(groupId)
	err := repositories.GetUsersFromGroup(intGroupId, resGroupusers, &groupusers)
	if err != nil {
		return err
	}
	if len(groupusers) > 0 {
		for _, groupuser := range groupusers {
			user_id := groupuser.UserID
			user := models.UserOnlyWithID{UserID: user_id}
			resGroupusers.Users = append(resGroupusers.Users, user)
		}
	}
	return nil
}

func CheckGroupExistsById(groupId string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupExistsById(groupId, &exists)
	if err != nil {
		return false, errors.New("group not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckGroupExistsByKey(key string, orgId string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupExistsByKey(key, orgId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func AddUsersToGroup(groupId string, users models.AddUsersToGroup) error {
	for _, user := range users.Users {
		userId := fmt.Sprint(user.UserID)
		exists, err := CheckUserAlreadyInGroup(groupId, userId)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		groupuser := models.GroupUser{}
		intGroupId, _ := strconv.Atoi(groupId)
		intUserId, _ := strconv.Atoi(userId)
		groupuser.GroupID = intGroupId
		groupuser.UserID = intUserId
		err = repositories.AddUserToGroup(groupuser)
		if err != nil {
			return err
		}

	}
	return nil

}

func CheckUserAlreadyInGroup(groupId string, userId string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupAlreadyInGroup(groupId, userId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
