package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResourceActions(c *gin.Context) {
	resourceActions := []models.ResourceAction{}
	res_id := string(c.Param("res_id"))
	exists, err := repositories.CheckResourceExistsById(res_id)
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
	repositories.GetResourceActions(&resourceActions, res_id)
	c.JSON(http.StatusOK, &resourceActions)
}

func GetResourceAction(c *gin.Context) {
	var resourceAction models.ResourceAction
	res_id := string(c.Param("res_id"))
	resact_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resact_exists, resact_err := repositories.CheckResourceActionExistsById(resact_id)
	if resact_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resact_exists {
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	repositories.GetResourceAction(&resourceAction, resact_id)
	c.JSON(http.StatusOK, &resourceAction)

}

func CreateResourceAction(c *gin.Context) {
	var resourceAction models.ResourceAction
	res_id := string(c.Param("res_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	if err := c.ShouldBindJSON(&resourceAction); err != nil {
		if resourceAction.Key == "" || len(resourceAction.Key) < 4 || resourceAction.Name == "" || len(resourceAction.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	int_res_id, _ := strconv.Atoi(res_id)
	resourceAction.ResourceID = int_res_id
	exists, err := repositories.CheckResourceActionExistsByKey(resourceAction.Key, res_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Resource Action already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action already exists"})
		return
	}
	repositories.CreateResourceActionAction(&resourceAction)
	c.JSON(http.StatusOK, &resourceAction)

}

func DeleteResourceAction(c *gin.Context) {
	var resourceAction models.ResourceAction
	res_id := string(c.Param("res_id"))
	resact_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resact_exists, resact_err := repositories.CheckResourceActionExistsById(resact_id)
	if resact_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resact_exists {
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	repositories.DeleteResourceAction(&resourceAction, resact_id)
	c.JSON(http.StatusOK, "")
}

func UpdateResourceAction(c *gin.Context) {
	var resourceAction models.ResourceAction
	var reqResourceAction models.ResourceAction
	res_id := string(c.Param("res_id"))
	resact_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resact_exists, resact_err := repositories.CheckResourceActionExistsById(resact_id)
	if resact_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resact_exists {
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	repositories.UpdateResourceAction(&resourceAction, &reqResourceAction, resact_id)
	c.JSON(http.StatusOK, &resourceAction)
}
