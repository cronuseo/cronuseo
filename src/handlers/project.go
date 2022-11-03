package handlers

import (
	"errors"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetProjects(projects *[]models.Project, orgId string) {
	repositories.GetProjects(projects, orgId)
}

func GetProject(project *models.Project, projId string) {
	repositories.GetProject(project, projId)
}

func CreateProject(project *models.Project) {
	repositories.CreateProject(project)
}

func DeleteProject(project *models.Project, projId string) {
	repositories.DeleteAllResources(projId)
	repositories.DeleteProject(project, projId)
}

func UpdateProject(project *models.Project, reqProject *models.Project, projId string) {
	repositories.GetProject(project, projId)
	project.Name = reqProject.Name
	project.Description = reqProject.Description
	repositories.UpdateProject(project)
}

func CheckProjectExistsById(projId string) (bool, error) {
	var exists bool
	err := repositories.CheckProjectExistsById(projId, exists)
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
	err := repositories.CheckProjectExistsByKey(key, orgId, exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
