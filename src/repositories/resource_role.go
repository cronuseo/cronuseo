package repositories

import (
	"log"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resource_id string, projects *[]models.ResourceRole) error {
	err := config.DB.Select(projects, "SELECT * FROM resource_role WHERE resource_id = $1", resource_id)
	if err != nil {
		return err
	}
	return nil
}

func GetResourceRole(resource_id string, id string, resource_role *models.ResourceRole) error {
	users := []models.UserID{}
	groups := []models.GroupID{}
	err := config.DB.Get(resource_role,
		"SELECT * FROM resource_role WHERE resource_id = $1 AND resource_role_id = $2", resource_id, id)
	if err != nil {
		return err
	}
	err = config.DB.Select(&users, "SELECT user_id FROM user_resource_role WHERE resource_role_id = $1", id)
	if err != nil {
		return err
	}
	err = config.DB.Select(&groups, "SELECT group_id FROM group_resource_role WHERE resource_role_id = $1", id)
	if err != nil {
		return err
	}
	resource_role.Users = users
	resource_role.Groups = groups
	return nil
}

func CreateResourceRole(resource_id string, resource_role *models.ResourceRole) error {

	var resource_role_id string

	tx, err := config.DB.Begin()

	if err != nil {
		return err
	}
	// add resource role
	{
		stmt, err := tx.Prepare(
			`INSERT INTO resource_role(resource_role_key,name,resource_id) VALUES($1, $2, $3) RETURNING resource_role_id`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		err = stmt.QueryRow(resource_role.Key, resource_role.Name, resource_id).Scan(&resource_role_id)

		if err != nil {
			return err
		}
	}

	// add users to esource role
	if len(resource_role.Users) > 0 {
		stmt, err := tx.Prepare("INSERT INTO user_resource_role(resource_role_id,user_id) VALUES ($1,$2)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for _, user := range resource_role.Users {
			_, err = stmt.Exec(resource_role_id, user.UserID)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	// add groups to esource role
	if len(resource_role.Groups) > 0 {
		stmt, err := tx.Prepare("INSERT INTO group_resource_role(resource_role_id,group_id) VALUES ($1,$2)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for _, group := range resource_role.Groups {
			_, err = stmt.Exec(resource_role_id, group.GroupID)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}
	resource_role.ID = resource_role_id
	return nil
}

func DeleteResourceRole(resource_id string, id string) error {

	tx, err := config.DB.Begin()

	if err != nil {
		return err
	}
	// remove all users from resource role
	{
		stmt, err := tx.Prepare(`DELETE FROM user_resource_role WHERE resource_role_id = $1`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(id)

		if err != nil {
			return err
		}
	}

	// remove all groups from resource role
	{
		stmt, err := tx.Prepare(`DELETE FROM group_resource_role WHERE esource_role_id = $1`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(id)

		if err != nil {
			return err
		}
	}

	// remove resource role
	{
		stmt, err := tx.Prepare(`DELETE FROM resource_role WHERE resource_id = $1 AND resource_role_id = $2`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(resource_id, id)
		if err != nil {
			return err
		}
	}

	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateResourceRole(resource_role *models.ResourceRole) error {
	stmt, err := config.DB.Prepare("UPDATE resource_role SET name = $1 WHERE resource_id = $2 AND resource_role_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource_role.Name, resource_role.ResourceID, resource_role.ID)
	if err != nil {
		return err
	}
	return nil
}

func PatchResourceRole(resource_role_id string, resourecRolePatch *models.ResourceRolePatchRequest) error {
	tx, err := config.DB.Begin()

	if err != nil {
		return err
	}

	for _, operation := range resourecRolePatch.Operations {
		switch operation.Path {

		case "users":
			{
				switch operation.Operation {
				case "add":
					if len(operation.Values) > 0 {
						stmt, err := tx.Prepare("INSERT INTO user_resource_role(resource_role_id,user_id) VALUES ($1,$2)")
						if err != nil {
							log.Fatal(err)
						}
						defer stmt.Close()

						for _, userId := range operation.Values {
							var exists bool
							user_error := IsUserInRole(resource_role_id, userId.Value, &exists)
							if exists || user_error != nil {
								continue
							}
							_, err = stmt.Exec(resource_role_id, userId.Value)
							if err != nil {
								log.Fatal(err)
							}

						}
					}
				case "remove":
					if len(operation.Values) > 0 {
						stmt, err := tx.Prepare("DELETE FROM user_resource_role WHERE resource_role_id = $1 AND user_id = $2")
						if err != nil {
							log.Fatal(err)
						}
						defer stmt.Close()

						for _, userId := range operation.Values {
							var exists bool
							user_error := IsUserInRole(resource_role_id, userId.Value, &exists)
							if !exists || user_error != nil {
								continue
							}
							_, err = stmt.Exec(resource_role_id, userId.Value)
							if err != nil {
								log.Fatal(err)
							}

						}
					}
				}

			}

		case "groups":
			{
				switch operation.Operation {
				case "add":
					if len(operation.Values) > 0 {
						stmt, err := tx.Prepare("INSERT INTO group_resource_role(resource_role_id,group_id) VALUES ($1,$2)")
						if err != nil {
							log.Fatal(err)
						}
						defer stmt.Close()

						for _, groupId := range operation.Values {
							var exists bool
							user_error := IsGroupInRole(resource_role_id, groupId.Value, &exists)
							if exists || user_error != nil {
								continue
							}
							_, err = stmt.Exec(resource_role_id, groupId.Value)
							if err != nil {
								log.Fatal(err)
							}

						}
					}
				case "remove":
					if len(operation.Values) > 0 {
						stmt, err := tx.Prepare("DELETE FROM group_resource_role WHERE resource_role_id = $1 AND group_id = $2")
						if err != nil {
							log.Fatal(err)
						}
						defer stmt.Close()

						for _, groupId := range operation.Values {
							var exists bool
							user_error := IsUserInRole(resource_role_id, groupId.Value, &exists)
							if !exists || user_error != nil {
								continue
							}
							_, err = stmt.Exec(resource_role_id, groupId.Value)
							if err != nil {
								log.Fatal(err)
							}

						}
					}
				}

			}
		}
	}

	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}

	return nil
}
func CheckResourceRoleExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT resource_role_id FROM resource_role WHERE resource_role_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckResourceRoleExistsByKey(resource_id string, key string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT resource_role_key FROM resource_role WHERE resource_id = $1 AND project_key = $2)",
		resource_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func IsUserInRole(resource_role_id string, userId string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT user_id FROM user_resource_role WHERE resource_role_id = $1 AND user_id = $2)",
		resource_role_id, userId).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func IsGroupInRole(resource_role_id string, groupId string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT user_id FROM group_resource_role WHERE resource_role_id = $1 AND group_id = $2)",
		resource_role_id, groupId).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
