package organization

import (
	"net/http"
	"testing"

	"github.com/shashimalcse/cronuseo/internal/middleware"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAPI(t *testing.T) {
	logger := test.InitLogger()
	router := test.MockRouter()
	repo := &mockRepository{orgs: []mongo_entity.Organization{
		{ID: primitive.NewObjectID(), Identifier: "test", DisplayName: "test"},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger))
	header := middleware.MockAuthHeader()

	tests := []test.APITestCase{
		{Name: "get all", Method: "GET", URL: "/organization", Body: "", Header: nil, WantStatus: http.StatusOK, WantResponse: ""},
		{Name: "create ok", Method: "POST", URL: "/organization", Body: `{"identifier":"test2", "display_name":"test2"}`, Header: header, WantStatus: http.StatusCreated, WantResponse: ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
