package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/:org_id/user")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type resource struct {
	service Service
}

// @Description Get user by ID.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Produce     json
// @Success     200 {object}  entity.User
// @failure     404,500
// @Router      /{org_id}/user/{id} [get]
func (r resource) get(c echo.Context) error {
	user, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, user)
}

// @Description Get all users.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param name query string false "name"
// @Param limit query integer false "limit"
// @Param cursor query integer false "cursor"
// @Produce     json
// @Success     200 {array}  entity.UserQueryResponse
// @failure     500
// @Router      /{org_id}/user [get]
func (r resource) query(c echo.Context) error {
	var filter Filter
	org_id := c.Param("org_id")
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	users, err := r.service.Query(c.Request().Context(), org_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	response := entity.UserQueryResponse{}
	maxUserID := -1
	minUserID := 10000
	for _, user := range users {
		newUser := entity.UserResult{ID: user.ID, Username: user.Username, FirstName: user.FirstName, LastName: user.LastName,
			OrgID: user.OrgID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
		newUser.Links = entity.UserLinks{Self: "/" + org_id + "/user/" + user.ID}
		response.Results = append(response.Results, newUser)
		if i, err := strconv.Atoi(user.LogicalKey); err == nil {
			if maxUserID < i {
				maxUserID = i
			}
			if minUserID > i {
				minUserID = i
			}
		}
	}
	response.Size = len(users)
	response.Limit = filter.Limit
	if len(users) > 0 {
		response.Cursor = maxUserID
		links := entity.Links{}
		links.Self = "/" + org_id + "/user/"
		if filter.Name != "" {
			links.Self += "?name=" + filter.Name
		}
		links.Self += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor)
		if len(users) == filter.Limit {
			links.Next = "/" + org_id + "/user/"
			if filter.Name != "" {
				links.Next += "?name=" + filter.Name
			}
			links.Next += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(response.Cursor)
		}
		if filter.Cursor != 0 {
			links.Prev = "/" + org_id + "/user/"
			if filter.Name != "" {
				links.Prev += "?name=" + filter.Name
			}
			links.Prev += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor-filter.Limit)
		}
		response.Links = links
	}
	return c.JSON(http.StatusOK, response)
}

// @Description Create user.
// @Tags        User
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateUserRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.User
// @failure     400,403,500
// @Router      /{org_id}/user [post]
func (r resource) create(c echo.Context) error {
	var input CreateUserRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	user, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, user)
}

// @Description Update user.
// @Tags        User
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.User
// @failure     400,403,404,500
// @Router      /{org_id}/user/{id} [put]
func (r resource) update(c echo.Context) error {
	var input UpdateUserRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	user, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, user)
}

// @Description Delete user.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/user/{id} [delete]
func (r resource) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
