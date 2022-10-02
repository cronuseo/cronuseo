package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetProjects(c *gin.Context) {
	projects := []models.Project{}
	checkOrganizationExists(c)
	config.DB.Model(&models.Project{}).Where("organization_id = ?", c.Param("org_id")).Find(&projects)
	c.JSON(http.StatusOK, &projects)
}

func CreateProjects(c *gin.Context) {
	var project models.Project
	checkOrganizationExists(c)
	if err := c.ShouldBindJSON(&project); err != nil {
		if project.Key == "" || len(project.Key) < 4 || project.Name == "" || len(project.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	org_id, _ := strconv.Atoi(c.Param("org_id"))
	project.OrganizationID = org_id
	exists, err := checkProjectExistsByKey(&project)
	fmt.Print(project)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project already exists"})
		return
	} else {
		config.DB.Create(&project)
		c.JSON(http.StatusOK, &project)
	}

}

func DeleteProjects(c *gin.Context) {
	var project models.Project
	checkOrganizationExists(c)
	exists, err := checkProjectExistsById(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).Delete(&project)
	c.JSON(http.StatusOK, "")
}

func UpdateProjects(c *gin.Context) {
	var project models.Project
	var reqProject models.Project
	checkOrganizationExists(c)
	if err := c.ShouldBindJSON(&reqProject); err != nil {
		if reqProject.Name == "" || len(reqProject.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
			return
		}
	}
	exists, err := checkProjectExistsById(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).First(&project)
	project.Name = reqProject.Name
	project.Description = reqProject.Description
	config.DB.Save(&project)
	c.JSON(http.StatusOK, &project)
}

/*
**

**
 */
func checkProjectExistsByKey(project *models.Project) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("key = ?", project.Key).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkProjectExistsById(c *gin.Context) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("id = ?", c.Param("id")).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkOrganizationExists(c *gin.Context) {
	var exists bool
	err := config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("id = ?", c.Param("org_id")).Find(&exists).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
}
