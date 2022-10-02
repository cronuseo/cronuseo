package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/utils"
	"golang.org/x/oauth2"
)

func HandleAuthentication(c *gin.Context) {
	url := utils.GithubOauthConfig.AuthCodeURL(utils.GetRandomState())
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallback(c *gin.Context) {
	if c.Request.FormValue("stage") != utils.RandomState {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	token, err := utils.GithubOauthConfig.Exchange(oauth2.NoContext, c.Request.FormValue("code"))

	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.JSON(http.StatusOK, token)
}
