package handlers

import (
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
