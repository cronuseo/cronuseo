package user

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (User, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]User, error)
	Create(ctx context.Context, org_id string, input CreateUserRequest) (User, error)
	Update(ctx context.Context, org_id string, id string, input UpdateUserRequest) (User, error)
	Delete(ctx context.Context, org_id string, id string) (User, error)
}

type User struct {
	entity.User
}

type CreateUserRequest struct {
	Username  string          `json:"username" db:"username"`
	FirstName string          `json:"firstname" db:"firstname"`
	LastName  string          `json:"lastname" db:"lastname"`
	OrgID     string          `json:"-" db:"org_id"`
	Roles     []entity.RoleID `json:"roles"`
}

func (m CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required),
	)
}

type UpdateUserRequest struct {
	FirstName string `json:"firstname" db:"firstname"`
	LastName  string `json:"lastname" db:"lastname"`
}

func (m UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&m)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) Get(ctx context.Context, org_id string, id string) (User, error) {
	user, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		return User{}, &util.NotFoundError{Path: "User"}
	}
	return User{user}, nil
}

func (s service) Create(ctx context.Context, org_id string, req CreateUserRequest) (User, error) {

	if err := req.Validate(); err != nil {
		return User{}, &util.InvalidInputError{}
	}

	//check organization exists
	exists, _ := s.repo.ExistByKey(ctx, req.Username)
	if exists {
		return User{}, &util.AlreadyExistsError{Path: "User"}
	}

	id := entity.GenerateID()
	err := s.repo.Create(ctx, org_id, entity.User{
		ID:        id,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Roles:     req.Roles,
	})
	if err != nil {
		return User{}, err
	}

	return s.Get(ctx, org_id, id)
}

func (s service) Update(ctx context.Context, org_id string, id string, req UpdateUserRequest) (User, error) {
	if err := req.Validate(); err != nil {
		return User{}, &util.InvalidInputError{}
	}

	user, err := s.Get(ctx, org_id, id)
	if err != nil {
		return user, &util.NotFoundError{Path: "User"}
	}
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	if err := s.repo.Update(ctx, org_id, user.User); err != nil {
		return user, err
	}
	return user, nil
}

func (s service) Delete(ctx context.Context, org_id string, id string) (User, error) {
	user, err := s.Get(ctx, org_id, id)
	if err != nil {
		return User{}, &util.NotFoundError{Path: "User"}
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		return User{}, err
	}
	return user, nil
}

type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]User, error) {
	items, err := s.repo.Query(ctx, org_id, filter)
	if err != nil {
		return nil, err
	}
	result := []User{}
	for _, item := range items {
		result = append(result, User{item})
	}
	return result, nil
}
