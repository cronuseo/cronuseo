package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/shashimalcse/Cronuseo/utils"
	"golang.org/x/oauth2"
)

func HandleAuthentication(c echo.Context) error {
	url := utils.GithubOauthConfig.AuthCodeURL(utils.GetRandomState())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallback(c echo.Context) error {
	if c.FormValue("stage") != utils.RandomState {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	token, err := utils.GithubOauthConfig.Exchange(oauth2.NoContext, c.FormValue("code"))

	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	return c.JSON(http.StatusOK, token)
}
