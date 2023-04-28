package organization

import (
	"context"
	"fmt"
	"testing"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/test"
	"github.com/shashimalcse/cronuseo/internal/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_service(t *testing.T) {
	logger := test.Init()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// Create new organization.

	_, err := s.Create(ctx, CreateOrganizationRequest{
		Identifier:  "test",
		DisplayName: "test",
	})
	assert.Nil(t, err)
	// assert.NotEmpty(t, id)

	// assert.NotEmpty(t, org.ID)
}

type mockRepository struct {
	orgs []mongo_entity.Organization
}

func (m mockRepository) Get(ctx context.Context, id string) (*mongo_entity.Organization, error) {
	fmt.Println(m.orgs)
	for _, org := range m.orgs {
		fmt.Println("checking")
		if org.ID.Hex() == id {
			return &org, nil
		}
	}
	return nil, &util.NotFoundError{Path: "Organization"}
}
func (m mockRepository) Create(ctx context.Context, organization mongo_entity.Organization) (string, error) {
	id := primitive.NewObjectID()
	organization.ID = id
	m.orgs = append(m.orgs, organization)
	return id.Hex(), nil
}
func (m mockRepository) Query(ctx context.Context) ([]mongo_entity.Organization, error) {
	return m.orgs, nil
}
func (m mockRepository) Delete(ctx context.Context, id string) error {
	return nil
}
func (m mockRepository) RefreshAPIKey(ctx context.Context, apiKey string, id string) error {
	return nil
}
func (m mockRepository) CheckOrgExistById(ctx context.Context, id string) (bool, error) {
	return true, nil
}
func (m mockRepository) CheckOrgExistByIdentifier(ctx context.Context, identifier string) (bool, error) {
	return false, nil
}
