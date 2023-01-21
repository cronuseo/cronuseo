package permission

import (
	"context"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/auth"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := permission{service}
	router := r.Group("/:org_id/permission")
	config := echojwt.Config{
		SigningKey: []byte(auth.SecretKey),
	}
	router.Use(echojwt.WithConfig(config))
	router.POST("/check_actions", res.checkActions)
	router.PATCH("/update", res.patchPermissions)
}

type permission struct {
	service Service
}

// @Description Delete tuple.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/delete [post]
func (r permission) delete(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckTuple(context.Background(), c.Param("org"), "permission", input, false)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org_id path string true "Organization"
// @Param request body CheckActionsRequest true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org_id}/permission/check_actions [post]
func (r permission) checkActions(c echo.Context) error {
	var input CheckActionsRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow := r.service.CheckActions(context.Background(), c.Param("org_id"), "permission", input)

	return c.JSON(http.StatusOK, allow)
}

// @Description Patch Permissions.
// @Tags        Permission
// @Accept      json
// @Param org_id path string true "Organization"
// @Param request body PermissionPatchRequest true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org_id}/permission/update [patch]
func (r permission) patchPermissions(c echo.Context) error {
	var input PermissionPatchRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	_ = r.service.PatchPermissions(context.Background(), c.Param("org_id"), "permission", input)

	return c.JSON(http.StatusOK, "")
}
