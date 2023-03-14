package monitoring

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"go.uber.org/zap"
)

type Service interface {
	GetAllowedData(ctx context.Context, org_id string) (entity.AllowedData, error)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return service{repo: repo, logger: logger}
}

// Get organization by id.
func (s service) GetAllowedData(ctx context.Context, org_id string) (entity.AllowedData, error) {

	data, err := s.repo.GetAllowed(ctx, org_id)

	return data, err
}
