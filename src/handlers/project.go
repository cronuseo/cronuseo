package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetProjects(projects *[]models.Project, orgId string) error {
	return repositories.GetProjects(projects, orgId)
}

func GetProject(project *models.Project, projId string) error {
	return repositories.GetProject(project, projId)
}

func CreateProject(project *models.Project) error {
	return repositories.CreateProject(project)
}

func DeleteProject(project *models.Project, projId string) error {
	err := repositories.DeleteAllResources(projId)
	if err != nil {
		return err
	}
	return repositories.DeleteProject(project, projId)
}

func UpdateProject(project *models.Project, reqProject *models.Project, projId string) error {
	err := repositories.GetProject(project, projId)
	if err != nil {
		return err
	}
	project.Name = reqProject.Name
	project.Description = reqProject.Description
	return repositories.UpdateProject(project)
}

func CheckProjectExistsById(projId string) (bool, error) {
	var exists bool
	err := repositories.CheckProjectExistsById(projId, &exists)
	if err != nil {
		return false, errors.New("project not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckProjectExistsByKey(key string, orgId string) (bool, error) {
	var exists bool
	err := repositories.CheckProjectExistsByKey(key, orgId, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
