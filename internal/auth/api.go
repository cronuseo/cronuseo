package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service, orgService organization.Service) {
	res := admin{service: service, orgService: orgService}
	r.POST("/register", res.register)
	r.POST("/login", res.login)
	r.POST("/logout", res.logout)
	config := echojwt.Config{
		SigningKey: []byte(SecretKey),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("/me", res.getMe)
}

type admin struct {
	service    Service
	orgService organization.Service
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
	token, err := r.service.Login(c.Request().Context(), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, TokenResponse{Token: token})
}

// @Description Logout.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Success     200
// @failure     400,403,500
// @Router      /logout [post]
func (r admin) logout(c echo.Context) error {
	cookie, err := r.service.Logout(c.Request().Context())
	if err != nil {
		return util.HandleError(err)
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "Success")
}

// @Description GetMe.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Success     200
// @failure     400,403,500
// @Router      /me [get]
func (r admin) getMe(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["sub"].(string)
	adminuser, err := r.service.GetMe(c.Request().Context(), name)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusOK, adminuser)
}
