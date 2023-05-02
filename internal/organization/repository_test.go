package organization

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/shashimalcse/cronuseo/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {

	db := test.DB(t)
	repo := NewRepository(db)

	ctx := context.Background()

	// Check default organization exist.
	bool, err := repo.CheckOrgExistByIdentifier(ctx, "super")
	assert.Nil(t, err)
	assert.Equal(t, true, bool)

	// Get all organizations.
	orgs, err := repo.Query(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(orgs))

	// Default organization.
	defaultOrg := orgs[0]

	// Generate new API key.
	key := make([]byte, 32)
	_, err = rand.Read(key)
	assert.Nil(t, err)

	newAPIKey := base64.StdEncoding.EncodeToString(key)

	// Refresh API key.
	err = repo.RefreshAPIKey(ctx, newAPIKey, defaultOrg.ID.Hex())
	assert.Nil(t, err)
}
