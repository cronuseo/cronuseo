package repositories

import (
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
	err := config.DB.Get(resource_role, "SELECT * FROM resource_role WHERE resource_id = $1 AND resource_role_id = $2", resource_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateResourceRole(resource_id string, resource_role *models.ResourceRole) error {

	stmt, err := config.DB.Prepare(
		"INSERT INTO resource_role(resource_role_key,name,resource_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource_role.Key, resource_role.Name, resource_id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteResourceRole(resource_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM resource_role WHERE resource_id = $1 AND resource_role_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource_id, id)
	if err != nil {
		return err
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
