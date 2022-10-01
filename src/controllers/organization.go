package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetOrganization(c *gin.Context) {
	orgs := []models.Organization{}
	config.DB.Find(&orgs)
	c.JSON(http.StatusOK, &orgs)
}

func CreateOrganization(c *gin.Context) {
	var orgs models.Organization
	c.BindJSON(&orgs)
	count := int64(0)
	err := config.DB.Model(&models.Organization{}).Where("key = ?", orgs.Key).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	exists := count > 0
	if exists {
		c.JSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization already exists"})
	} else {
		c.BindJSON(&orgs)
		config.DB.Create(&orgs)
		c.JSON(http.StatusOK, &orgs)
	}

}

func DeleteOrganization(c *gin.Context) {
	var orgs models.Organization
	config.DB.Where("id = ?", c.Param("id")).Delete(&orgs)
	c.JSON(http.StatusOK, &orgs)
}

func UpdateOrganization(c *gin.Context) {
	var orgs models.Organization
	config.DB.Where("id = ?", c.Param("id")).First(&orgs)
	c.BindJSON(&orgs)
	config.DB.Save(&orgs)
	c.JSON(http.StatusOK, &orgs)
}
