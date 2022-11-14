package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetGroups(tenantId string, groups *[]models.Group) error {
	return repositories.GetGroups(tenantId, groups)
}

func GetGroup(tenantId string, id string, group *models.Group) error {
	return repositories.GetGroup(tenantId, id, group)
}

func CreateGroup(tenantId string, group *models.Group) error {
	group.Users = RemoveDuplicateUsers(group.Users)
	group.Users = GetExistsUsers(tenantId, group.Users)
	err := repositories.CreateGroup(tenantId, group)
	return err

}

func DeleteGroup(tenantId string, id string) error {
	return repositories.DeleteGroup(tenantId, id)
}

func UpdateGroup(tenantId string, id string, group *models.Group, reqGroup *models.GroupUpdateRequest) error {
	err := repositories.GetGroup(tenantId, id, group)
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

func CheckGroupExistsByKey(tenantId string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupExistsByKey(tenantId, key, &exists)
	return exists, err
}

func PatchGroup(tenantId string, groupId string, group *models.Group, groupPatch *models.GroupPatchRequest) error {
	for _, operation := range groupPatch.Operations {
		operation.Users = RemoveDuplicateUsers(operation.Users)
	}
	err := repositories.PatchGroup(groupId, groupPatch)
	if err != nil {
		return err
	}
	return GetGroup(tenantId, groupId, group)

}

func GetExistsUsers(tenantId string, users []models.UserID) []models.UserID {
	existsUsers := []models.UserID{}
	for _, user := range users {
		exists, _ := CheckUserExistsByTenant(tenantId, user.UserID)
		if exists {
			existsUsers = append(existsUsers, user)
		}

	}
	return existsUsers
}

func RemoveDuplicateUsers(userIDs []models.UserID) []models.UserID {

	processed := make(map[models.UserID]struct{})

	uniqUserIDs := make([]models.UserID, 0)
	for _, uid := range userIDs {
		if _, ok := processed[uid]; ok {
			continue
		}
		uniqUserIDs = append(uniqUserIDs, uid)
		processed[uid] = struct{}{}
	}

	return uniqUserIDs
}
