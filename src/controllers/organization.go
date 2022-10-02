package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetOrganizations(c *gin.Context) {
	orgs := []models.Organization{}
	config.DB.Find(&orgs)
	c.JSON(http.StatusOK, &orgs)
}

func CreateOrganization(c *gin.Context) {
	var orgs models.Organization
	if err := c.ShouldBindJSON(&orgs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
		return
	}
	exists, err := checkOrganizationExistsByKey(&orgs)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Organization already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization already exists"})
		return
	}
	config.DB.Create(&orgs)
	c.JSON(http.StatusOK, &orgs)
}

func DeleteOrganization(c *gin.Context) {
	var orgs models.Organization
	exists, err := checkOrganizationExistsById(c)
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
	config.DB.Where("id = ?", c.Param("id")).Delete(&orgs)
	c.JSON(http.StatusOK, "")
}

func UpdateOrganization(c *gin.Context) {
	var orgs models.Organization
	var reqOrgs models.Organization
	if err := c.ShouldBindJSON(&reqOrgs); err != nil {
		if reqOrgs.Name == "" || len(reqOrgs.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
			return
		}
	}
	exists, err := checkOrganizationExistsById(c)
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
	config.DB.Where("id = ?", c.Param("id")).First(&orgs)
	orgs.Name = reqOrgs.Name
	config.DB.Save(&orgs)
	c.JSON(http.StatusOK, &orgs)
}

/*
**

**
 */
func checkOrganizationExistsByKey(orgs *models.Organization) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("key = ?", orgs.Key).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func checkOrganizationExistsById(c *gin.Context) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("id = ?", c.Param("id")).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
