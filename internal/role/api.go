package role

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := role{service}
	router := r.Group("/:org_id/role")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
	router.GET("/user/:user_id", res.QueryByUserID)
}

type role struct {
	service Service
}

// @Description Get role by ID.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Param id path string true "Role ID"
// @Produce     json
// @Success     200 {object}  entity.Role
// @failure     404,500
// @Router      /{org_id}/role/{id} [get]
func (r role) get(c echo.Context) error {
	role, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, role)
}

// @Description Get all roles.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Param name query string false "name"
// @Param limit query integer false "limit"
// @Param cursor query integer false "cursor"
// @Produce     json
// @Success     200 {array}  entity.Role
// @failure     500
// @Router      /{org_id}/role [get]
func (r role) query(c echo.Context) error {
	var filter Filter
	org_id := c.Param("org_id")
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	roles, err := r.service.Query(c.Request().Context(), org_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	response := entity.RoleQueryResponse{}
	maxRoleID := -1
	minRoleID := 10000
	for _, role := range roles {
		newAction := entity.RoleResult{ID: role.ID, Name: role.Name, Key: role.Key,
			OrgID: role.OrgID, CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt}
		newAction.Links = entity.RoleLinks{Self: "/" + role.OrgID + "/role/" + role.ID}
		response.Results = append(response.Results, newAction)
		if i, err := strconv.Atoi(role.LogicalKey); err == nil {
			if maxRoleID < i {
				maxRoleID = i
			}
			if minRoleID > i {
				minRoleID = i
			}
		}

	}
	response.Size = len(roles)
	response.Limit = filter.Limit
	if len(roles) > 0 {
		response.Cursor = maxRoleID
		links := entity.Links{}
		links.Self = "/" + org_id + "/role/"
		if filter.Name != "" {
			links.Self += "?name=" + filter.Name
		}
		links.Self += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor)
		if len(roles) == filter.Limit {
			links.Next = "/" + org_id + "/role/"
			if filter.Name != "" {
				links.Next += "?name=" + filter.Name
			}
			links.Next += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(response.Cursor)
		}
		if filter.Cursor != 0 {
			links.Prev = "/" + org_id + "/role/"
			if filter.Name != "" {
				links.Prev += "?name=" + filter.Name
			}
			links.Prev += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor-filter.Limit)
		}
		response.Links = links
	}
	return c.JSON(http.StatusOK, response)
}

// @Description Create role.
// @Tags        Role
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateRoleRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Role
// @failure     400,403,500
// @Router      /{org_id}/role [post]
func (r role) create(c echo.Context) error {
	var input CreateRoleRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	role, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, role)
}

// @Description Update role.
// @Tags        Role
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Role ID"
// @Param request body UpdateRoleRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Role
// @failure     400,403,404,500
// @Router      /{org_id}/role/{id} [put]
func (r role) update(c echo.Context) error {
	var input UpdateRoleRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	role, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, role)
}

// @Description Delete role.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Param id path string true "Role ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/role/{id} [delete]
func (r role) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Get all roles.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Produce     json
// @Success     200 {array}  entity.Role
// @failure     500
// @Router      /{org_id}/role/user/{user_id} [get]
func (r role) QueryByUserID(c echo.Context) error {
	var filter Filter
	org_id := c.Param("org_id")
	user_id := c.Param("user_id")
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	roles, err := r.service.QueryByUserID(c.Request().Context(), org_id, user_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	response := entity.RoleQueryResponse{}
	maxRoleID := -1
	minRoleID := 10000
	for _, role := range roles {
		newAction := entity.RoleResult{ID: role.ID, Name: role.Name, Key: role.Key,
			OrgID: role.OrgID, CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt}
		newAction.Links = entity.RoleLinks{Self: "/" + role.OrgID + "/role/" + role.ID}
		response.Results = append(response.Results, newAction)
		if i, err := strconv.Atoi(role.LogicalKey); err == nil {
			if maxRoleID < i {
				maxRoleID = i
			}
			if minRoleID > i {
				minRoleID = i
			}
		}

	}
	response.Size = len(roles)
	response.Limit = filter.Limit
	if len(roles) > 0 {
		response.Cursor = maxRoleID
		links := entity.Links{}
		links.Self = "/" + org_id + "/role/"
		if filter.Name != "" {
			links.Self += "?name=" + filter.Name
		}
		links.Self += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor)
		if len(roles) == filter.Limit {
			links.Next = "/" + org_id + "/role/"
			if filter.Name != "" {
				links.Next += "?name=" + filter.Name
			}
			links.Next += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(response.Cursor)
		}
		if filter.Cursor != 0 {
			links.Prev = "/" + org_id + "/role/"
			if filter.Name != "" {
				links.Prev += "?name=" + filter.Name
			}
			links.Prev += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor-filter.Limit)
		}
		response.Links = links
	}
	return c.JSON(http.StatusOK, response)
}
