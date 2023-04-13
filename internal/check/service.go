package check

import (
	"github.com/shashimalcse/cronuseo/internal/entity"
	"go.uber.org/zap"
)

type Service interface {
}

type Tuple struct {
	entity.Tuple
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}
