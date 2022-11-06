package handlers

import (
	"errors"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetUsers(users *[]models.User, orgId string) {
	repositories.GetUsers(users, orgId)
}

func GetUser(user *models.UserWithGroup, userId string) {
	repositories.GetUsersWithGroups(userId, user)
}

func CreateUser(user *models.User) {
	repositories.CreateUser(user)
}

func DeleteUser(user *models.User, userId string) {
	repositories.DeleteAllResources(userId)
	repositories.DeleteUser(user, userId)
}

func UpdateUser(user *models.User, reqUser *models.User, userId string) {
	repositories.GetUser(user, userId)
	user.FirstName = reqUser.FirstName
	user.LastName = reqUser.LastName
	repositories.UpdateUser(user)
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
