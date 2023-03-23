package resource

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/auth"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/:org_id/resource")
	config := echojwt.Config{
		SigningKey: []byte(auth.SecretKey),
	}
	router.Use(echojwt.WithConfig(config))
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	// router.DELETE("/:id", res.delete)
	// router.PUT("/:id", res.update)
}

type resource struct {
	service Service
}

// @Description Get resource by ID.
// @Tags        Resource
// @Param org_id path string true "Organization ID"
// @Param id path string true "Resource ID"
// @Produce     json
// @Success     200 {object}  entity.Resource
// @failure     404,500
// @Router      /{org_id}/resource/{id} [get]
func (r resource) get(c echo.Context) error {

	resource, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, resource)
}

// @Description Get all resources.
// @Tags        Resource
// @Param org_id path string true "Organization ID"
// @Param name query string false "name"
// @Param limit query integer false "limit"
// @Param cursor query integer false "cursor"
// @Produce     json
// @Success     200 {array}  entity.ResourceQueryResponse
// @failure     500
// @Router      /{org_id}/resource [get]
func (r resource) query(c echo.Context) error {

	org_id := c.Param("org_id")
	var filter Filter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	// Get all resources.
	resources, err := r.service.Query(c.Request().Context(), org_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	// response := entity.ResourceQueryResponse{}
	// maxResourceID := -1
	// minResourceID := 10000

	// // Create resource results for the response.
	// for _, resource := range resources {
	// 	newResource := entity.ResourceResult{ID: resource.ID, Name: resource.Name, Key: resource.Key,
	// 		OrgID: resource.OrgID, CreatedAt: resource.CreatedAt, UpdatedAt: resource.UpdatedAt}
	// 	newResource.Links = entity.ResourceLinks{Self: "/" + org_id + "/resource/" + resource.ID}
	// 	response.Results = append(response.Results, newResource)
	// 	if i, err := strconv.Atoi(resource.LogicalKey); err == nil {
	// 		if maxResourceID < i {
	// 			maxResourceID = i
	// 		}
	// 		if minResourceID > i {
	// 			minResourceID = i
	// 		}
	// 	}

	// }
	// // Pagination
	// response.Size = len(resources)
	// response.Limit = filter.Limit
	// if len(resources) > 0 {
	// 	response.Cursor = maxResourceID
	// 	links := entity.Links{}
	// 	links.Self = "/" + org_id + "/resource/"
	// 	if filter.Name != "" {
	// 		links.Self += "?name=" + filter.Name
	// 	}
	// 	links.Self += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor)
	// 	if len(resources) == filter.Limit {
	// 		links.Next = "/" + org_id + "/resource/"
	// 		if filter.Name != "" {
	// 			links.Next += "?name=" + filter.Name
	// 		}
	// 		links.Next += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(response.Cursor)
	// 	}
	// 	if filter.Cursor != 0 {
	// 		links.Prev = "/" + org_id + "/resource/"
	// 		if filter.Name != "" {
	// 			links.Prev += "?name=" + filter.Name
	// 		}
	// 		links.Prev += "&limit=" + strconv.Itoa(filter.Limit) + "&cursor=" + strconv.Itoa(filter.Cursor-filter.Limit)
	// 	}
	// 	response.Links = links
	// }
	return c.JSON(http.StatusOK, resources)
}

// @Description Create resource.
// @Tags        Resource
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateResourceRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Resource
// @failure     400,403,500
// @Router      /{org_id}/resource [post]
func (r resource) create(c echo.Context) error {

	var input CreateResourceRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	resource, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, resource)
}

// @Description Update resource.
// @Tags        Resource
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Resource ID"
// @Param request body UpdateResourceRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Resource
// @failure     400,403,404,500
// @Router      /{org_id}/resource/{id} [put]
// func (r resource) update(c echo.Context) error {

// 	var input UpdateResourceRequest
// 	if err := c.Bind(&input); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
// 	}

// 	resource, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
// 	if err != nil {
// 		return util.HandleError(err)
// 	}
// 	return c.JSON(http.StatusCreated, resource)
// }

// @Description Delete resource.
// @Tags        Resource
// @Param org_id path string true "Organization ID"
// @Param id path string true "Resource ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/resource/{id} [delete]
// func (r resource) delete(c echo.Context) error {

// 	_, err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
// 	if err != nil {
// 		return util.HandleError(err)
// 	}
// 	return c.JSON(http.StatusNoContent, "")
// }
