package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetUsers(users *[]models.User, org_id string) {
	config.DB.Model(&models.User{}).Where("organization_id = ?", org_id).Find(&users)
}

func GetUser(user *models.User, user_id string) {
	config.DB.Where("id = ?", user_id).First(&user)
}

func CreateUser(user *models.User) {
	config.DB.Create(&user)
}

func DeleteUser(user *models.User, user_id string) {
	DeleteAllResources(user_id)
	config.DB.Where("id = ?", user_id).Delete(&user)
}

func UpdateUser(user *models.User, reqUser *models.User, user_id string) {
	config.DB.Where("id = ?", user_id).First(&user)
	user.Name = reqUser.Name
	config.DB.Save(&user)
}

func CheckUserExistsById(user_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.User{}).Select("count(*) > 0").Where("id = ?", user_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("user not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckUserExistsByUsername(username string, org_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.User{}).Select("count(*) > 0").Where("key = ? AND organization_id = ?", key, org_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
