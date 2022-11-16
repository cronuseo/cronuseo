package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetGroups(orgId string, groups *[]models.Group) error {
	return repositories.GetGroups(orgId, groups)
}

func GetGroup(orgId string, id string, group *models.Group) error {
	return repositories.GetGroup(orgId, id, group)
}

func CreateGroup(orgId string, group *models.Group) error {
	group.Users = RemoveDuplicateUsers(group.Users)
	group.Users = GetExistsUsers(orgId, group.Users)
	err := repositories.CreateGroup(orgId, group)
	return err

}

func DeleteGroup(orgId string, id string) error {
	return repositories.DeleteGroup(orgId, id)
}

func UpdateGroup(orgId string, id string, group *models.Group, reqGroup *models.GroupUpdateRequest) error {
	err := repositories.GetGroup(orgId, id, group)
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

func CheckGroupExistsByKey(orgId string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckGroupExistsByKey(orgId, key, &exists)
	return exists, err
}

func PatchGroup(orgId string, groupId string, group *models.Group, groupPatch *models.GroupPatchRequest) error {
	for _, operation := range groupPatch.Operations {
		operation.Users = RemoveDuplicateUsers(operation.Users)
	}
	err := repositories.PatchGroup(groupId, groupPatch)
	if err != nil {
		return err
	}
	return GetGroup(orgId, groupId, group)

}

func GetExistsUsers(orgId string, users []models.UserID) []models.UserID {
	existsUsers := []models.UserID{}
	for _, user := range users {
		exists, _ := CheckUserExistsByTenant(orgId, user.UserID)
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
