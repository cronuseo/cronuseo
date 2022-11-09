package repositories

import (
	"strconv"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetUsers(users *[]models.User, org_id string) error {
	return config.DB.Model(&models.User{}).Where("organization_id = ?", org_id).Find(&users).Error
}

func GetUser(user *models.User, userId string) error {
	return config.DB.Where("id = ?", userId).First(&user).Error
}

func CreateUser(user *models.User) error {
	return config.DB.Create(&user).Error
}

func DeleteUser(user *models.User, userId string) error {
	return config.DB.Where("id = ?", userId).Delete(&user).Error
}

func GetUsersWithGroups(userId string, resUser *models.UserWithGroup,
	groupusers *[]models.GroupUser) error {
	intUserId, _ := strconv.Atoi(userId)
	err := config.DB.Model(&models.User{}).Select("id", "username", "name", "organization_id").Where(
		"id = ?", userId).Find(&resUser).Error
	if err != nil {
		return err
	}
	return config.DB.Model(&models.GroupUser{}).Where("user_id = ?", intUserId).Find(&groupusers).Error

}

func UpdateUser(user *models.User) error {
	return config.DB.Save(&user).Error
}

func CheckUserExistsById(userId string, exists *bool) error {
	return config.DB.Model(&models.User{}).Select("count(*) > 0").Where("id = ?",
		userId).Find(exists).Error
}

func CheckUserExistsByUsername(username string, orgId string, exists *bool) error {
	return config.DB.Model(&models.User{}).Select("count(*) > 0").Where(
		"username = ? AND organization_id = ?", username, orgId).Find(exists).Error
}
