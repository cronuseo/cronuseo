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

func GetResourceActions(c *gin.Context) {
	resourceActions := []models.ResourceAction{}
	checkProjectExists(c)
	config.DB.Model(&models.Resource{}).Where("resource_id = ?", c.Param("res_id")).Find(&resourceActions)
	c.JSON(http.StatusOK, &resourceActions)
}

func CreateResourceAction(c *gin.Context) {
	var resourceAction models.ResourceAction
	resExists := checkResourceExists(c)
	if !resExists {
		return
	}
	if err := c.ShouldBindJSON(&resourceAction); err != nil {
		if resourceAction.Key == "" || len(resourceAction.Key) < 4 || resourceAction.Name == "" || len(resourceAction.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	org_id, _ := strconv.Atoi(c.Param("res_id"))
	resourceAction.ResourceID = org_id
	exists, err := checkResourceActionExistsByKey(&resourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Resource Action already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action already exists"})
		return
	} else {
		config.DB.Create(&resourceAction)
		c.JSON(http.StatusOK, &resourceAction)
	}

}

func DeleteResourceAction(c *gin.Context) {
	var resourceAction models.ResourceAction
	resExists := checkResourceExists(c)
	if !resExists {
		return
	}
	exists, err := checkResourceActionExistsById(c)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).Delete(&resourceAction)
	c.JSON(http.StatusOK, "")
}

func UpdateResourceAction(c *gin.Context) {
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
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	config.DB.Where("id = ?", c.Param("id")).First(&resourceAction)
	resourceAction.Name = reqResourceAction.Name
	resourceAction.Description = reqResourceAction.Description
	config.DB.Save(&resourceAction)
	c.JSON(http.StatusOK, &resourceAction)
}

func checkResourceActionExistsByKey(resourceAction *models.ResourceAction) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceAction{}).Select("count(*) > 0").Where("key = ?", resourceAction.Key).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkResourceActionExistsById(c *gin.Context) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceAction{}).Select("count(*) > 0").Where("id = ?", c.Param("id")).Find(&exists).Error
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

func checkResourceExists(c *gin.Context) bool {
	var exists bool
	err := config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("id = ?", c.Param("res_id")).Find(&exists).Error
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return false
	}
	if !exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return false
	}

	return true
}
