package repositories

import (
	"errors"
	"strconv"

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

func GetUsersWithGroups(user_id string, resUser *models.UserWithGroup) {
	groupusers := []models.GroupUser{}
	int_user_id, _ := strconv.Atoi(user_id)
	config.DB.Model(&models.User{}).Select("id", "username", "name", "organization_id").Where("id = ?", user_id).Find(&resUser)
	config.DB.Model(&models.GroupUser{}).Where("user_id = ?", int_user_id).Find(&groupusers)
	if len(groupusers) > 0 {
		for _, groupuser := range groupusers {
			group_id := groupuser.UserID
			user := models.GroupOnlyWithID{GroupID: group_id}
			resUser.Groups = append(resUser.Groups, user)
		}
	}

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
	err := config.DB.Model(&models.User{}).Select("count(*) > 0").Where("username = ? AND organization_id = ?", username, org_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
