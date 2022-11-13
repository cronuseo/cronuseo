package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetGroups(tenant_id string, groups *[]models.Group) error {
	return repositories.GetGroups(tenant_id, groups)
}

func GetGroup(tenant_id string, id string, group *models.Group) error {
	return repositories.GetGroup(tenant_id, id, group)
}

func CreateGroup(tenant_id string, group *models.Group) error {
	return repositories.CreateGroup(tenant_id, group)
}

func DeleteGroup(tenant_id string, id string) error {
	return repositories.DeleteGroup(tenant_id, id)
}

func UpdateGroup(tenant_id string, id string, group *models.Group, reqGroup *models.GroupUpdateRequest) error {
	err := repositories.GetGroup(tenant_id, id, group)
	if err != nil {
		return err
	}
	group.Name = reqGroup.Name
	return repositories.UpdateGroup(group)
}

func CheckGroupExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupExistsById(id, &exists)
	return exists, err
}

func CheckGroupExistsByKey(tenant_id string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupExistsByKey(tenant_id, key, &exists)
	return exists, err
}

// func AddUserToGroup(groupId string, userId string) error {
// 	groupuser := models.GroupUser{}
// 	intGroupId, _ := strconv.Atoi(groupId)
// 	intUserId, _ := strconv.Atoi(userId)
// 	groupuser.GroupID = intGroupId
// 	groupuser.UserID = intUserId
// 	return repositories.AddUserToGroup(groupuser)

// }

// func GetUsersFromGroup(groupId string, resGroupusers *models.GroupUsers) error {
// 	var groupusers []models.GroupUser
// 	intGroupId, _ := strconv.Atoi(groupId)
// 	err := repositories.GetUsersFromGroup(intGroupId, resGroupusers, &groupusers)
// 	if err != nil {
// 		return err
// 	}
// 	if len(groupusers) > 0 {
// 		for _, groupuser := range groupusers {
// 			user_id := groupuser.UserID
// 			user := models.UserOnlyWithID{UserID: user_id}
// 			resGroupusers.Users = append(resGroupusers.Users, user)
// 		}
// 	}
// 	return nil
// }

// func CheckGroupExistsById(groupId string) (bool, error) {
// 	var exists bool
// 	err := repositories.CheckGroupExistsById(groupId, &exists)
// 	if err != nil {
// 		return false, errors.New("group not exists")
// 	}
// 	if exists {
// 		return true, nil
// 	} else {
// 		return false, nil
// 	}
// }

// func CheckGroupExistsByKey(key string, orgId string) (bool, error) {
// 	var exists bool
// 	err := repositories.CheckGroupExistsByKey(key, orgId, &exists)
// 	if err != nil {
// 		return false, errors.New("")
// 	}
// 	if exists {
// 		return true, nil
// 	} else {
// 		return false, nil
// 	}
// }

// func AddUsersToGroup(groupId string, users models.AddUsersToGroup) error {
// 	for _, user := range users.Users {
// 		userId := fmt.Sprint(user.UserID)
// 		exists, err := CheckUserAlreadyInGroup(groupId, userId)
// 		if err != nil {
// 			return err
// 		}
// 		if exists {
// 			continue
// 		}
// 		groupuser := models.GroupUser{}
// 		intGroupId, _ := strconv.Atoi(groupId)
// 		intUserId, _ := strconv.Atoi(userId)
// 		groupuser.GroupID = intGroupId
// 		groupuser.UserID = intUserId
// 		err = repositories.AddUserToGroup(groupuser)
// 		if err != nil {
// 			return err
// 		}

// 	}
// 	return nil

// }

// func CheckUserAlreadyInGroup(groupId string, userId string) (bool, error) {
// 	var exists bool
// 	err := repositories.CheckGroupAlreadyInGroup(groupId, userId, &exists)
// 	if err != nil {
// 		return false, errors.New("")
// 	}
// 	if exists {
// 		return true, nil
// 	} else {
// 		return false, nil
// 	}
// }
