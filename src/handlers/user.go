package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetUsers(orgId string, users *[]models.User) error {
	return repositories.GetUsers(orgId, users)
}

// func GetUser(resUser *models.UserWithGroup, userId string) error {
// 	var groupusers []models.GroupUser
// 	err := repositories.GetUsersWithGroups(userId, resUser, &groupusers)
// 	if err != nil {
// 		return err
// 	}
// 	if len(groupusers) > 0 {
// 		for _, groupuser := range groupusers {
// 			groupId := groupuser.UserID
// 			user := models.GroupOnlyWithID{GroupID: groupId}
// 			resUser.Groups = append(resUser.Groups, user)
// 		}
// 	}
// 	return nil
// }

func GetUser(orgId string, userId string, user *models.User) error {

	err := repositories.GetUser(orgId, userId, user)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(orgId string, user *models.User) error {
	return repositories.CreateUser(orgId, user)
}

func DeleteUser(orgId string, userId string) error {
	return repositories.DeleteUser(orgId, userId)
}

func UpdateUser(orgId string, userId string, user *models.User, reqUser *models.UserUpdateRequest) error {
	err := repositories.GetUser(orgId, userId, user)
	if err != nil {
		return err
	}
	user.FirstName = reqUser.FirstName
	user.LastName = reqUser.LastName
	return repositories.UpdateUser(user)
}

func CheckUserExistsById(userId string) (bool, error) {
	var exists bool
	err := repositories.CheckUserExistsById(userId, &exists)
	return exists, err
}

func CheckUserExistsByTenant(orgId string, userId string) (bool, error) {
	var exists bool
	err := repositories.CheckUserExistsByTenant(orgId, userId, &exists)
	return exists, err
}

func CheckUserExistsByUsername(orgId string, username string) (bool, error) {
	var exists bool
	err := repositories.CheckUserExistsByUsername(orgId, username, &exists)
	return exists, err
}
