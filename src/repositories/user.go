package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetUsers(tenant_id string, users *[]models.User) error {
	err := config.DB.Select(users, "SELECT * FROM tenant_user WHERE tenant_id = $1", tenant_id)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(tenant_id string, id string, user *models.User) error {
	err := config.DB.Get(user, "SELECT * FROM tenant_user WHERE tenant_id = $1 AND user_id = $2", tenant_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(tenant_id string, user *models.User) error {
	stmt, err := config.DB.Prepare(
		"INSERT INTO tenant_user(username,first_name,last_name,tenant_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.FirstName, user.LastName, tenant_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(tenant_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM tenant_user WHERE tenant_id = $1 AND user_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tenant_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *models.User) error {
	stmt, err := config.DB.Prepare(
		"UPDATE tenant_user SET first_name = $1, last_name = $2, tenant_id = $3 WHERE user_id = $4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.FirstName, user.LastName, user.TenantID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT user_id FROM tenant_user WHERE user_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserExistsByTenant(tenant_id string, id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT user_id FROM tenant_user WHERE tenant_id = $1 AND user_id = $2)",
		tenant_id, id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserExistsByUsername(tenant_id string, username string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT username FROM tenant_user WHERE tenant_id = $1 AND username = $2)",
		tenant_id, username).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
