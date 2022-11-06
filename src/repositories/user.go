package repositories

import (
	"strconv"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetUsers(users *[]models.User, org_id string) {
	config.DB.Model(&models.User{}).Where("organization_id = ?", org_id).Find(&users)
}

func GetUser(user *models.User, userId string) {
	config.DB.Where("id = ?", userId).First(&user)
}

func CreateUser(user *models.User) {
	config.DB.Create(&user)
}

func DeleteUser(user *models.User, userId string) {
	config.DB.Where("id = ?", userId).Delete(&user)
}

func GetUsersWithGroups(userId string, resUser *models.UserWithGroup) {
	var groupusers []models.GroupUser
	intUserId, _ := strconv.Atoi(userId)
	config.DB.Model(&models.User{}).Select("id", "username", "name", "organization_id").Where(
		"id = ?", userId).Find(&resUser)
	config.DB.Model(&models.GroupUser{}).Where("user_id = ?", intUserId).Find(&groupusers)
	if len(groupusers) > 0 {
		for _, groupuser := range groupusers {
			groupId := groupuser.UserID
			user := models.GroupOnlyWithID{GroupID: groupId}
			resUser.Groups = append(resUser.Groups, user)
		}
	}

}

func UpdateUser(user *models.User) {
	config.DB.Save(&user)
}

func CheckUserExistsById(userId string, exists *bool) error {
	return config.DB.Model(&models.User{}).Select("count(*) > 0").Where("id = ?",
		userId).Find(exists).Error
}

func CheckUserExistsByUsername(username string, orgId string, exists *bool) error {
	return config.DB.Model(&models.User{}).Select("count(*) > 0").Where(
		"username = ? AND organization_id = ?", username, orgId).Find(exists).Error
}
