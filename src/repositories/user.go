package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetUsers(org_id string, users *[]models.User) error {
	err := config.DB.Select(users, "SELECT * FROM organization_user WHERE org_id = $1", org_id)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(org_id string, id string, user *models.User) error {
	err := config.DB.Get(user, "SELECT * FROM organization_user WHERE org_id = $1 AND user_id = $2", org_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(org_id string, user *models.User) error {
	stmt, err := config.DB.Prepare(
		"INSERT INTO organization_user(username,first_name,last_name,org_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.FirstName, user.LastName, org_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(org_id string, id string) error {

	tx, err := config.DB.Begin()

	if err != nil {
		return err
	}

	// remove user from group
	{
		stmt, err := tx.Prepare(`DELETE FROM group_user WHERE user_id = $1`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(id)

		if err != nil {
			return err
		}
	}

	// remove user from resource role
	{
		stmt, err := tx.Prepare(`DELETE FROM user_resource_role WHERE user_id = $1`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(id)

		if err != nil {
			return err
		}
	}

	// remove user
	{
		stmt, err := tx.Prepare(`DELETE FROM organization_user WHERE org_id = $1 AND user_id = $2`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(org_id, id)
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

func UpdateUser(user *models.User) error {
	stmt, err := config.DB.Prepare(
		"UPDATE organization_user SET first_name = $1, last_name = $2, org_id = $3 WHERE user_id = $4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.FirstName, user.LastName, user.OrgID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT user_id FROM organization_user WHERE user_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserExistsByTenant(org_id string, id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT user_id FROM organization_user WHERE org_id = $1 AND user_id = $2)",
		org_id, id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserExistsByUsername(org_id string, username string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT username FROM organization_user WHERE org_id = $1 AND username = $2)",
		org_id, username).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
