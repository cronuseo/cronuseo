package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetProjects(tenant_id string, projects *[]models.Project) error {
	err := config.DB.Select(projects, "SELECT * FROM project WHERE tenant_id = $1", tenant_id)
	if err != nil {
		return err
	}
	return nil
}

func GetProject(tenant_id string, id string, project *models.Project) error {
	err := config.DB.Get(project, "SELECT * FROM project WHERE tenant_id = $1 AND project_id = $2", tenant_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateProject(tenant_id string, project *models.Project) error {
	stmt, err := config.DB.Prepare("INSERT INTO project(project_key,name,tenant_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project.Key, project.Name, tenant_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProject(tenant_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM project WHERE tenant_id = $1 AND project_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tenant_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProject(project *models.Project) error {
	stmt, err := config.DB.Prepare("UPDATE project SET name = $1 WHERE tenant_id = $2 AND project_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project.Name, project.TenantID, project.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckProjectExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT project_id FROM project WHERE project_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckProjectExistsByKey(tenant_id string, key string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT project_key FROM project WHERE tenant_id = $1 AND project_key = $2)",
		tenant_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
