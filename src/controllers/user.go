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

func GetUsers(c *gin.Context) {
	users := []models.User{}
	org_id := string(c.Param("org_id"))
	exists, err := repositories.CheckOrganizationExistsById(org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	repositories.GetUsers(&users, org_id)
	c.JSON(http.StatusOK, &users)
}

func GetUser(c *gin.Context) {
	var user models.UserWithGroup
	org_id := string(c.Param("org_id"))
	user_id := string(c.Param("id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	user_exists, user_err := repositories.CheckUserExistsById(user_id)
	if user_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !user_exists {
		config.Log.Info("User not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User not exists"})
		return
	}
	repositories.GetUsersWithGroups(user_id, &user)
	c.JSON(http.StatusOK, &user)

}

func CreateUser(c *gin.Context) {
	var user models.User
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		if user.Username == "" || len(user.Username) < 4 || user.Name == "" || len(user.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	int_org_id, _ := strconv.Atoi(org_id)
	user.OrganizationID = int_org_id
	exists, err := repositories.CheckUserExistsByUsername(user.Username, org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("User already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User already exists"})
		return
	}
	repositories.CreateUser(&user)
	c.JSON(http.StatusOK, &user)

}

func DeleteUser(c *gin.Context) {
	var user models.User
	user_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	user_exists, user_err := repositories.CheckUserExistsById(user_id)
	if user_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !user_exists {
		config.Log.Info("User not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User not exists"})
		return
	}
	repositories.DeleteUser(&user, user_id)
	c.JSON(http.StatusOK, "")
}

func UpdateUser(c *gin.Context) {
	var user models.User
	var reqUser models.User
	user_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	user_exists, user_err := repositories.CheckUserExistsById(user_id)
	if user_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !user_exists {
		config.Log.Info("User not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User not exists"})
		return
	}
	repositories.UpdateUser(&user, &reqUser, user_id)
	c.JSON(http.StatusOK, &user)
}
