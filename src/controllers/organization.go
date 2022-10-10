package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetOrganizations(c *gin.Context) {
	orgs := []models.Organization{}
	repositories.GetOrganizations(&orgs)
	c.JSON(http.StatusOK, &orgs)
}

func GetOrganization(c *gin.Context) {
	var org models.Organization
	org_id := string(c.Param("id"))
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
	repositories.GetOrganization(&org, org_id)
	c.JSON(http.StatusOK, &org)

}

func CreateOrganization(c *gin.Context) {
	var org models.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
		return
	}
	exists, err := repositories.CheckOrganizationExistsByKey(org.Key)
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
	repositories.CreateOrganization(&org)
	c.JSON(http.StatusOK, &org)
}

func DeleteOrganization(c *gin.Context) {
	var org models.Organization
	org_id := string(c.Param("id"))
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
	repositories.DeleteOrganization(&org, org_id)
	c.JSON(http.StatusOK, "")
}

func UpdateOrganization(c *gin.Context) {
	var org models.Organization
	org_id := string(c.Param("id"))
	var reqOrg models.Organization
	if err := c.ShouldBindJSON(&reqOrg); err != nil {
		if reqOrg.Name == "" || len(reqOrg.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Invalid inputs. Please check your inputs"})
			return
		}
	}
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
	repositories.UpdateOrganization(&org, &reqOrg, org_id)
	c.JSON(http.StatusOK, &org)
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
