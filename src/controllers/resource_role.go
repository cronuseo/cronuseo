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

func GetResourceRoles(c *gin.Context) {
	resourceRoles := []models.ResourceRole{}
	checkProjectExists(c)
	config.DB.Model(&models.Resource{}).Where("resource_id = ?", c.Param("res_id")).Find(&resourceRoles)
	c.JSON(http.StatusOK, &resourceRoles)
}

func GetResourceRole(c *gin.Context) {
	var resRole models.ResourceRole
	exists, err := checkResourceRoleExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource Role Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).First(&resRole)
	c.JSON(http.StatusOK, &resRole)

}

func CreateResourceRole(c *gin.Context) {
	var resourceRole models.ResourceRole
	resExists := checkResourceExists(c)
	if !resExists {
		return
	}
	if err := c.ShouldBindJSON(&resourceRole); err != nil {
		if resourceRole.Key == "" || len(resourceRole.Key) < 4 || resourceRole.Name == "" || len(resourceRole.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	org_id, _ := strconv.Atoi(c.Param("res_id"))
	resourceRole.ResourceID = org_id
	exists, err := checkResourceRoleExistsByKey(&resourceRole)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Resource Role already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role already exists"})
		return
	} else {
		config.DB.Create(&resourceRole)
		c.JSON(http.StatusOK, &resourceRole)
	}

}

func DeleteResourceRole(c *gin.Context) {
	var resourceRole models.ResourceRole
	resExists := checkResourceExists(c)
	if !resExists {
		return
	}
	exists, err := checkResourceRoleExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).Delete(&resourceRole)
	c.JSON(http.StatusOK, "")
}

func UpdateResourceRole(c *gin.Context) {
	var resourceAction models.ResourceAction
	var reqResourceAction models.ResourceAction
	resExists := checkProjectExists(c)
	if !resExists {
		return
	}
	if err := c.ShouldBindJSON(&reqResourceAction); err != nil {
		if reqResourceAction.Name == "" || len(reqResourceAction.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
			return
		}
	}
	exists, err := checkResourceActionExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).First(&resourceAction)
	resourceAction.Name = reqResourceAction.Name
	resourceAction.Description = reqResourceAction.Description
	config.DB.Save(&resourceAction)
	c.JSON(http.StatusOK, &resourceAction)
}

func checkResourceRoleExistsByKey(resourceRole *models.ResourceRole) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("key = ?", resourceRole.Key).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkResourceRoleExistsById(c *gin.Context) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("id = ?", c.Param("id")).Find(&exists).Error
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
