package integration

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	orgJson = `{
    "key" : "shsashimal1",
    "name" : "shashimal senarath"}`
)

func TestCreateOrganization(t *testing.T) {
	// Setup
	println("rocket")
	connectDB()
	e := echo.New()
	StartServer(e)
	req := httptest.NewRequest(http.MethodPost, "/organization", strings.NewReader(orgJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	println(rec.Body)
	// Assertions
	if assert.NoError(t, controllers.CreateOrganization(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, orgJson, rec.Body.String())
	}
}
