package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetProjects(tenant_id string, projects *[]models.Project) error {
	return repositories.GetProjects(tenant_id, projects)
}

func GetProject(tenant_id string, id string, project *models.Project) error {
	return repositories.GetProject(tenant_id, id, project)
}

func CreateProject(tenant_id string, project *models.Project) error {
	return repositories.CreateProject(tenant_id, project)
}

func DeleteProject(tenant_id string, id string) error {
	return repositories.DeleteProject(tenant_id, id)
}

func UpdateProject(tenant_id string, id string, project *models.Project,
	reqProject *models.ProjectUpdateRequest) error {
	err := repositories.GetProject(tenant_id, id, project)
	if err != nil {
		return err
	}
	project.Name = reqProject.Name
	return repositories.UpdateProject(project)
}

func CheckProjectExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckProjectExistsById(id, &exists)
	return exists, err
}

func CheckProjectExistsByKey(tenant_id string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckProjectExistsByKey(tenant_id, key, &exists)
	return exists, err
}
