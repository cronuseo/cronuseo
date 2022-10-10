package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResources(c *gin.Context) {
	resources := []models.Resource{}
	checkProjectExists(c)
	config.DB.Model(&models.Resource{}).Where("project_id = ?", c.Param("proj_id")).Find(&resources)
	c.JSON(http.StatusOK, &resources)
}

func GetResource(c *gin.Context) {
	var res models.Resource
	exists, err := checkResourceExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).First(&res)
	c.JSON(http.StatusOK, &res)

}

func CreateResource(c *gin.Context) {
	var resource models.Resource
	projExists := checkProjectExists(c)
	if !projExists {
		return
	}
	if err := c.ShouldBindJSON(&resource); err != nil {
		if resource.Key == "" || len(resource.Key) < 4 || resource.Name == "" || len(resource.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	org_id, _ := strconv.Atoi(c.Param("proj_id"))
	resource.ProjectID = org_id
	exists, err := checkResourceExistsByKey(&resource)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Resource already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource already exists"})
		return
	} else {
		config.DB.Create(&resource)
		c.JSON(http.StatusOK, &resource)
	}

}

func DeleteResource(c *gin.Context) {
	var resource models.Resource
	projExists := checkProjectExists(c)
	if !projExists {
		return
	}
	exists, err := checkResourceExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).Delete(&resource)
	c.JSON(http.StatusOK, "")
}

func UpdateResource(c *gin.Context) {
	var resource models.Resource
	var reqResource models.Resource
	projExists := checkProjectExists(c)
	if !projExists {
		return
	}
	if err := c.ShouldBindJSON(&reqResource); err != nil {
		if reqResource.Name == "" || len(reqResource.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
			return
		}
	}
	exists, err := checkProjectExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).First(&resource)
	resource.Name = reqResource.Name
	resource.Description = reqResource.Description
	config.DB.Save(&resource)
	c.JSON(http.StatusOK, &resource)
}

func checkResourceExistsByKey(resource *models.Resource) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("key = ?", resource.Key).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkResourceExistsById(c *gin.Context) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("id = ?", c.Param("id")).Find(&exists).Error
	if err != nil {
		config.Log.Panic("Server Error!")
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkProjectExists(c *gin.Context) bool {
	var exists bool
	err := config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("id = ?", c.Param("proj_id")).Find(&exists).Error
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return false
	}
	if !exists {
		config.Log.Info("Project not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
		return false
	}

	return true
}
