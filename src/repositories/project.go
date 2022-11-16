package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetProjects(org_id string, projects *[]models.Project) error {
	err := config.DB.Select(projects, "SELECT * FROM project WHERE org_id = $1", org_id)
	if err != nil {
		return err
	}
	return nil
}

func GetProject(org_id string, id string, project *models.Project) error {
	err := config.DB.Get(project, "SELECT * FROM project WHERE org_id = $1 AND project_id = $2", org_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateProject(org_id string, project *models.Project) error {
	stmt, err := config.DB.Prepare("INSERT INTO project(project_key,name,org_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project.Key, project.Name, org_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProject(org_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM project WHERE org_id = $1 AND project_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(org_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProject(project *models.Project) error {
	stmt, err := config.DB.Prepare("UPDATE project SET name = $1 WHERE org_id = $2 AND project_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(project.Name, project.OrgID, project.ID)
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

func CheckProjectExistsByKey(org_id string, key string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT project_key FROM project WHERE org_id = $1 AND project_key = $2)",
		org_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
