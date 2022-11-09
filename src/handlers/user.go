package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetUsers(users *[]models.User, orgId string) error {
	return repositories.GetUsers(users, orgId)
}

func GetUser(resUser *models.UserWithGroup, userId string) error {
	var groupusers []models.GroupUser
	err := repositories.GetUsersWithGroups(userId, resUser, &groupusers)
	if err != nil {
		return err
	}
	if len(groupusers) > 0 {
		for _, groupuser := range groupusers {
			groupId := groupuser.UserID
			user := models.GroupOnlyWithID{GroupID: groupId}
			resUser.Groups = append(resUser.Groups, user)
		}
	}
	return nil
}

func CreateUser(user *models.User) error {
	return repositories.CreateUser(user)
}

func DeleteUser(user *models.User, userId string) error {
	return repositories.DeleteUser(user, userId)
}

func UpdateUser(user *models.User, reqUser *models.User, userId string) error {
	err := repositories.GetUser(user, userId)
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
	if err != nil {
		return false, errors.New("user not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckUserExistsByUsername(username string, orgId string) (bool, error) {
	var exists bool
	err := repositories.CheckUserExistsByUsername(username, orgId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
