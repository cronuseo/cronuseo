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
	Login(ctx context.Context, req AdminUserRequest) (string, error)
	Logout(ctx context.Context) (*http.Cookie, error)
	GetMe(ctx context.Context, user_id string) (entity.AdminUser, error)
}

type AdminUser struct {
	entity.AdminUser
}

type TokenResponse struct {
	Token string `json:"token"`
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
		IsSuper:  true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s service) Login(ctx context.Context, req AdminUserRequest) (string, error) {

	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", &util.NotFoundError{Path: "Username"}
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return "", &util.InvalidInputError{Path: "Password"}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		return "", &util.SystemError{}
	}

	return token, nil

}

func (s service) Logout(ctx context.Context) (*http.Cookie, error) {

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)

	return cookie, nil
}

func (s service) GetMe(ctx context.Context, user_id string) (entity.AdminUser, error) {

	user, err := s.repo.GetUserByID(ctx, user_id)
	if err != nil {
		return entity.AdminUser{}, &util.NotFoundError{Path: "User id"}
	}
	return user, nil

}
