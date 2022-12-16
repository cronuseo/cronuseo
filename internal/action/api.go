package action

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := action{service}
	router := r.Group("/:resource_id/action")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type action struct {
	service Service
}

// @Description Get action by ID.
// @Tags        Action
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Action ID"
// @Produce     json
// @Success     200 {object}  entity.Action
// @failure     404,500
// @Router      /{resource_id}/action/{id} [get]
func (r action) get(c echo.Context) error {
	action, err := r.service.Get(c.Request().Context(), c.Param("resource_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, action)
}

// @Description Get all resources.
// @Tags        Action
// @Param resource_id path string true "Resource ID"
// @Param name query string false "name"
// @Param limit query integer false "limit"
// @Param cursor query integer false "cursor"
// @Produce     json
// @Success     200 {array}  entity.Action
// @failure     500
// @Router      /{resource_id}/action [get]
func (r action) query(c echo.Context) error {

	var filter Filter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	resource_id := c.Param("resource_id")
	actions, err := r.service.Query(c.Request().Context(), resource_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	response := entity.ActionQueryResponse{}
	maxActionID := -1
	minActionID := 10000
	for _, action := range actions {
		newAction := entity.ActionResult{ID: action.ID, Name: action.Name, Key: action.Key,
			ResourceID: action.ResourceID, CreatedAt: action.CreatedAt, UpdatedAt: action.UpdatedAt}
		newAction.Links = entity.ActionLinks{Self: "/" + action.ResourceID + "/action/" + action.ID}
		response.Results = append(response.Results, newAction)
		if i, err := strconv.Atoi(action.LogicalKey); err == nil {
			if maxActionID < i {
				maxActionID = i
			}
			if minActionID > i {
				minActionID = i
			}
		}

	}
	response.Size = len(actions)
	response.Limit = filter.Limit
	if len(actions) > 0 {
		response.Cursor = maxActionID
		links := entity.Links{}
		links.Self = "/" + resource_id + "/action/"
		if filter.Name != "" {
			links.Self += "?name=" + filter.Name
		}
		links.Self += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor)
		if len(actions) == filter.Limit {
			links.Next = "/" + resource_id + "/action/"
			if filter.Name != "" {
				links.Next += "?name=" + filter.Name
			}
			links.Next += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(response.Cursor)
		}
		if filter.Cursor != 0 {
			links.Prev = "/" + resource_id + "/action/"
			if filter.Name != "" {
				links.Prev += "?name=" + filter.Name
			}
			links.Prev += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor-filter.Limit)
		}
		response.Links = links
	}
	return c.JSON(http.StatusOK, response)
}

// @Description Create action.
// @Tags        Action
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param request body CreateActionRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Action
// @failure     400,403,500
// @Router      /{resource_id}/action [post]
func (r action) create(c echo.Context) error {
	var input CreateActionRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	action, err := r.service.Create(c.Request().Context(), c.Param("resource_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, action)
}

// @Description Update action.
// @Tags        Action
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Action ID"
// @Param request body UpdateActionRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Action
// @failure     400,403,404,500
// @Router      /{resource_id}/action/{id} [put]
func (r action) update(c echo.Context) error {
	var input UpdateActionRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	action, err := r.service.Update(c.Request().Context(), c.Param("resource_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, action)
}

// @Description Delete action.
// @Tags        Action
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Action ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{resource_id}/action/{id} [delete]
func (r action) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("resource_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
