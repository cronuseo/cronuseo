package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"golang.org/x/crypto/bcrypt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Register(ctx context.Context, adminUser AdminUserRequest) error
}

type AdminUser struct {
	entity.AdminUser
}

type AdminUserRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (m AdminUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Password, validation.Required),
		validation.Field(&m.Username, validation.Length(4, 30)),
		validation.Field(&m.Password, validation.Length(8, 30)),
	)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

const SecretKey = "secret"

func (s service) Register(ctx context.Context, req AdminUserRequest) error {

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

func (s service) Login(ctx context.Context, req AdminUserRequest) (*http.Cookie, error) {

	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, &util.NotFoundError{Path: "Username"}
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return nil, &util.InvalidInputError{Path: "Password"}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		return nil, &util.SystemError{}
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)

	return cookie, nil

}
