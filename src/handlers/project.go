package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetProjects(orgId string, projects *[]models.Project) error {
	return repositories.GetProjects(orgId, projects)
}

func GetProject(orgId string, id string, project *models.Project) error {
	return repositories.GetProject(orgId, id, project)
}

func CreateProject(orgId string, project *models.Project) error {
	return repositories.CreateProject(orgId, project)
}

func DeleteProject(orgId string, id string) error {
	return repositories.DeleteProject(orgId, id)
}

func UpdateProject(orgId string, id string, project *models.Project,
	reqProject *models.ProjectUpdateRequest) error {
	err := repositories.GetProject(orgId, id, project)
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

func CheckProjectExistsByKey(orgId string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckProjectExistsByKey(orgId, key, &exists)
	return exists, err
}
