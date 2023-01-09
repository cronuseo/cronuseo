package auth

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"golang.org/x/crypto/bcrypt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Register(ctx context.Context, adminUser RegisterAdminUser) error
}

type AdminUser struct {
	entity.AdminUser
}

type RegisterAdminUser struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (m RegisterAdminUser) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Password, validation.Required),
		validation.Field(&m.Username, validation.Length(4, 30)),
		validation.Field(&m.Password, validation.Length(8, 30)),
	)
}

type LoginRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) Register(ctx context.Context, req RegisterAdminUser) error {

	if err := req.Validate(); err != nil {
		return &util.InvalidInputError{}
	}

	exists, _ := s.repo.ExistByUsername(ctx, req.Username)
	if exists {
		return &util.InvalidInputError{}
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	id := entity.GenerateID()
	err := s.repo.Create(ctx, entity.AdminUser{
		ID:       id,
		Username: req.Username,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}
