package organization

import (
	"context"
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

	// successful creation
	org, err := s.Create(ctx, CreateOrganizationRequest{
		Identifier:  "test",
		DisplayName: "test",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, org.ID)
	assert.Equal(t, "test", org.Identifier)
	assert.Equal(t, "test", org.DisplayName)

	// validation error in creation
	_, err = s.Create(ctx, CreateOrganizationRequest{
		DisplayName: "test",
	})
	assert.NotNil(t, err)

}

type mockRepository struct {
	orgs []mongo_entity.Organization
}

func (m mockRepository) Get(ctx context.Context, id string) (*mongo_entity.Organization, error) {
	for _, org := range m.orgs {
		if org.ID.Hex() == id {
			return &org, nil
		}
	}
	return nil, &util.NotFoundError{Path: "Organization"}
}
func (m *mockRepository) Create(ctx context.Context, organization mongo_entity.Organization) (string, error) {
	id := primitive.NewObjectID()
	organization.ID = id
	m.orgs = append(m.orgs, organization)
	return id.Hex(), nil
}
func (m mockRepository) Query(ctx context.Context) ([]mongo_entity.Organization, error) {
	return m.orgs, nil
}
func (m mockRepository) Delete(ctx context.Context, id string) error {
	for i, org := range m.orgs {
		if org.ID.Hex() == id {
			m.orgs[i] = m.orgs[len(m.orgs)-1]
			m.orgs = m.orgs[:len(m.orgs)-1]
			break
		}
	}
	return nil
}
func (m *mockRepository) RefreshAPIKey(ctx context.Context, apiKey string, id string) error {
	for _, org := range m.orgs {
		if org.ID.Hex() == id {
			org.API_KEY = apiKey
		}
	}
	return &util.NotFoundError{Path: "Organization"}
}
func (m mockRepository) CheckOrgExistById(ctx context.Context, id string) (bool, error) {
	for _, org := range m.orgs {
		if org.ID.Hex() == id {
			return true, nil
		}
	}
	return false, nil
}
func (m mockRepository) CheckOrgExistByIdentifier(ctx context.Context, identifier string) (bool, error) {
	for _, org := range m.orgs {
		if org.Identifier == identifier {
			return true, nil
		}
	}
	return false, nil
}
