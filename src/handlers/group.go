package handlers

import (
	"errors"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
	"strconv"
)

func GetGroups(groups *[]models.Group, orgId string) {
	repositories.GetGroups(groups, orgId)
}

func GetGroup(group *models.Group, groupId string) {
	repositories.GetGroup(group, groupId)
}

func CreateGroup(group *models.Group) {
	repositories.CreateGroup(group)
}

func DeleteGroup(group *models.Group, groupId string) {
	repositories.DeleteAllResources(groupId)
	repositories.DeleteGroup(group, groupId)
}

func UpdateGroup(group *models.Group, reqGroup *models.Group, groupId string) {
	repositories.GetGroup(group, groupId)
	group.Name = reqGroup.Name
	group.Key = reqGroup.Key
	repositories.UpdateGroup(group)
}

func AddUserToGroup(groupId string, userId string) {
	groupuser := models.GroupUser{}
	intGroupId, _ := strconv.Atoi(groupId)
	intUserId, _ := strconv.Atoi(userId)
	groupuser.GroupID = intGroupId
	groupuser.UserID = intUserId
	repositories.AddUserToGroup(groupuser)

}

func GetUsersFromGroup(groupId string, resGroupusers *models.GroupUsers) {
	var groupusers []models.GroupUser
	intGroupId, _ := strconv.Atoi(groupId)
	repositories.GetUsersFromGroup(intGroupId, resGroupusers, groupusers)
	if len(groupusers) > 0 {
		for _, groupuser := range groupusers {
			user_id := groupuser.UserID
			user := models.UserOnlyWithID{UserID: user_id}
			resGroupusers.Users = append(resGroupusers.Users, user)
		}
	}

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
