package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceActions(resource_id string, resourceActions *[]models.ResourceAction) error {
	err := config.DB.Select(resourceActions, "SELECT * FROM resource_action WHERE resource_id = $1", resource_id)
	if err != nil {
		return err
	}
	return nil
}

func GetResourceAction(resource_id string, id string, resourceAction *models.ResourceAction) error {
	err := config.DB.Get(resourceAction,
		"SELECT * FROM resource_action WHERE resource_id = $1 AND resource_action_id = $2", resource_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateResourceAction(resource_id string, resourceAction *models.ResourceAction) error {
	stmt, err := config.DB.Prepare("INSERT INTO resource_action(resource_action_key,name,resource_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resourceAction.Key, resourceAction.Name, resource_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteResourceAction(resource_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM resource_action WHERE resource_id = $1 AND resource_action_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateResourceAction(resourceAction *models.ResourceAction) error {
	stmt, err := config.DB.Prepare(
		"UPDATE resource_action SET name = $1 WHERE resource_id = $2 AND resource_action_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resourceAction.Name, resourceAction.ResourceID, resourceAction.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckResourceActionExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT resource_action_id FROM resource_action WHERE resource_action_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckResourceActionExistsByKey(resource_id string, key string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT resource_action_key FROM resource_action WHERE resource_id = $1 AND resource_action_key = $2)",
		resource_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func GetActionKeyById(id string, key *string) error {
	err := config.DB.QueryRow("SELECT resource_action_key FROM resource_action WHERE resource_action_id = $1",
		id).Scan(key)
	if err != nil {
		return err
	}
	return nil
}
