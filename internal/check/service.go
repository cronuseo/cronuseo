package check

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"
)

type Service interface {
	Check(ctx context.Context, org string, req CheckRequest, apiKey string) (bool, error)
}

type CheckRequest struct {
	Username string `json:"username"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}

func (s service) Check(ctx context.Context, org string, req CheckRequest, apiKey string) (bool, error) {

	// Check resource already exists.
	validated, _ := s.repo.ValidateAPIKey(ctx, org, apiKey)
	if !validated {
		s.logger.Debug("API_KEY is not valid.")
		return false, &util.UnauthorizedError{}
	}
	return true, nil
}
