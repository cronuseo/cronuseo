package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := admin{service: service}
	router := r.Group("/auth")
	router.POST("/register", res.register)
	router.POST("/login", res.login)
	router.POST("/logout", res.logout)
	config := echojwt.Config{
		SigningKey: []byte(SecretKey),
	}
	router.Use(echojwt.WithConfig(config))
	router.GET("/me", res.getMe)
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
// @Router      /auth/register [post]
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
// @Router      /auth/login [post]
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
// @Router      /auth/logout [post]
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
// @Router      /auth/me [get]
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
