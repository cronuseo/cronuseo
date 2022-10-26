package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func Check(c *gin.Context) {
	var keys models.ResourceRoleToResourceActionKey
	if err := c.ShouldBindJSON(&keys); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
		return
	}

	allow, err := repositories.Check(keys.Resource, keys.ResourceRole, keys.ResourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if allow {
		c.JSON(http.StatusOK, "allowed")
	} else {
		c.JSON(http.StatusOK, "not allowed")
	}

}

func Check2(c echo.Context) error {
	var keys models.ResourceRoleToResourceActionKey
	if err := c.Bind(&keys); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: err.Error()})

	}
	if err := c.Validate(&keys); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: "Invalid inputs. Please check your inputs"})
	}
	allow, err := repositories.Check(keys.Resource, keys.ResourceRole, keys.ResourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if allow {
		return c.JSON(http.StatusOK, "allowed")
	} else {
		return c.JSON(http.StatusOK, "not allowed")
	}

}
