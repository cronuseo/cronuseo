package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResources(project_id string, projects *[]models.Resource) error {
	err := config.DB.Select(projects, "SELECT * FROM resource WHERE project_id = $1", project_id)
	if err != nil {
		return err
	}
	return nil
}

func GetResource(project_id string, id string, resource *models.Resource) error {
	err := config.DB.Get(resource, "SELECT * FROM resource WHERE project_id = $1 AND resource_id = $2", project_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateResource(project_id string, resource *models.Resource) error {

	stmt, err := config.DB.Prepare(
		"INSERT INTO resource(resource_key,name,project_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource.Key, resource.Name, project_id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteResource(project_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM resource WHERE project_id = $1 AND resource_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateResource(resource *models.Resource) error {
	stmt, err := config.DB.Prepare("UPDATE resource SET name = $1 WHERE project_id = $2 AND resource_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource.Name, resource.ProjectID, resource.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckResourceExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT resource_id FROM resource WHERE resource_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckResourceExistsByKey(project_id string, key string, exists *bool) error {
	err := config.DB.QueryRow(
		"SELECT exists (SELECT resource_key FROM resource WHERE project_id = $1 AND project_key = $2)",
		project_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func GetResourceKeyById(resourceId string, key *string) error {
	err := config.DB.QueryRow("SELECT resource_key FROM resource WHERE resource_id = $1", resourceId).Scan(key)
	if err != nil {
		return err
	}
	return nil
}

func GetProjectIDById(resourceId string, projectId *string) error {
	err := config.DB.QueryRow("SELECT project_id FROM resource WHERE resource_id = $1", resourceId).Scan(projectId)
	if err != nil {
		return err
	}
	return nil
}
