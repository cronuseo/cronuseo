package action

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, resource_id string, id string) (Action, error)
	Query(ctx context.Context, resource_id string, filter Filter) ([]Action, error)
	Create(ctx context.Context, resource_id string, input CreateActionRequest) (Action, error)
	Update(ctx context.Context, resource_id string, id string, input UpdateActionRequest) (Action, error)
	Delete(ctx context.Context, resource_id string, id string) (Action, error)
}

type Action struct {
	entity.Action
}

type CreateActionRequest struct {
	Key        string `json:"action_key" db:"action_key"`
	Name       string `json:"name" db:"name"`
	ResourceID string `json:"-" db:"resource_id"`
}

// Validating action creation request.
func (m CreateActionRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Key, validation.Required),
	)
}

type UpdateActionRequest struct {
	Name string `json:"name" db:"name"`
}

// Validating action update request.
func (m UpdateActionRequest) Validate() error {

	return validation.ValidateStruct(&m)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}

// Get action by id.
func (s service) Get(ctx context.Context, resource_id string, id string) (Action, error) {

	resource, err := s.repo.Get(ctx, resource_id, id)
	if err != nil {
		s.logger.Error("Error while getting the action.",
			zap.String("resource_id", resource_id),
			zap.String("action_id", id))
		return Action{}, &util.NotFoundError{Path: "Action"}
	}
	return Action{resource}, nil
}

// Create action.
func (s service) Create(ctx context.Context, resource_id string, req CreateActionRequest) (Action, error) {

	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating action creation request.")
		return Action{}, &util.InvalidInputError{}
	}

	exists, _ := s.repo.ExistByKey(ctx, resource_id, req.Key)
	if exists {
		s.logger.Debug("Action already exists.")
		return Action{}, &util.AlreadyExistsError{Path: "Action"}
	}

	// Generate unique id for action.
	id := entity.GenerateID()
	err := s.repo.Create(ctx, resource_id, entity.Action{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})

	if err != nil {
		s.logger.Error("Error while creating action.",
			zap.String("resource_id", resource_id),
			zap.String("action_key", req.Key))
		return Action{}, err
	}
	return s.Get(ctx, resource_id, id)
}

// Update action.
func (s service) Update(ctx context.Context, resource_id string, id string, req UpdateActionRequest) (Action, error) {

	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating action creation request.")
		return Action{}, &util.InvalidInputError{}
	}

	resource, err := s.Get(ctx, resource_id, id)
	if err != nil {
		s.logger.Debug("Action not exists.", zap.String("action_id", id))
		return Action{}, &util.NotFoundError{Path: "Action " + id + " not exists."}
	}
	resource.Name = req.Name
	if err := s.repo.Update(ctx, resource_id, resource.Action); err != nil {
		s.logger.Error("Error while updating action.",
			zap.String("resource_id", resource_id),
			zap.String("action_id", id))
		return resource, err
	}
	return resource, err
}

// Delete action.
func (s service) Delete(ctx context.Context, resource_id string, id string) (Action, error) {

	resource, err := s.Get(ctx, resource_id, id)
	if err != nil {
		s.logger.Error("Action not exists.", zap.String("action_id", id))
		return Action{}, &util.NotFoundError{Path: "Action"}
	}
	if err = s.repo.Delete(ctx, resource_id, id); err != nil {
		s.logger.Error("Error while deleting action.",
			zap.String("resource_id", resource_id),
			zap.String("action_id", id))
		return Action{}, err
	}
	return resource, err
}

// Pagination filter.
type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// Get all actions.
func (s service) Query(ctx context.Context, resource_id string, filter Filter) ([]Action, error) {

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	items, err := s.repo.Query(ctx, resource_id, filter)
	if err != nil {
		s.logger.Error("Error while retrieving all actions.",
			zap.String("resource_id", resource_id))
		return nil, err
	}
	result := []Action{}
	for _, item := range items {
		result = append(result, Action{item})
	}
	return result, nil
}
