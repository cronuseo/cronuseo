package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := admin{service}
	r.POST("/register", res.register)
	r.POST("/login", res.login)
}

type admin struct {
	service Service
}

// @Description Register.
// @Tags        Auth
// @Accept      json
// @Param request body AdminUserRequest true "body"
// @Produce     json
// @Success     200
// @failure     400,403,500
// @Router      /register [post]
func (r admin) register(c echo.Context) error {
	var input AdminUserRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	err := r.service.Register(c.Request().Context(), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, "")
}

// @Description Login.
// @Tags        Auth
// @Accept      json
// @Param request body AdminUserRequest true "body"
// @Produce     json
// @Success     200
// @failure     400,403,500
// @Router      /login [post]
func (r admin) login(c echo.Context) error {
	var input AdminUserRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	cookie, err := r.service.Login(c.Request().Context(), input)
	if err != nil {
		return util.HandleError(err)
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "Success")
}
