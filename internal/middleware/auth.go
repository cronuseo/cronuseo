package middleware

import (
	"fmt"
	"net/http"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shashimalcse/cronuseo/internal/config"
	"go.uber.org/zap"
)

func Auth(cfg *config.Config, logger *zap.Logger) echo.MiddlewareFunc {

	jwks, err := keyfunc.Get(cfg.Auth.JWKS, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			logger.Error("There was an error with the jwt.KeyFunc", zap.Error(err))
		},
	})

	if err != nil {
		logger.Error("Failed to create JWKs from resource at the given URL.", zap.Error(err))
	}

	// initialize JWT middleware instance
	return middleware.JWTWithConfig(middleware.JWTConfig{
		// skip public endpoints
		// Skipper: func(context echo.Context) bool {
		// 	return context.Path() == "/"
		// },
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			// a hack to convert jwt -> v4 tokens
			t, _, err := new(jwtv4.Parser).ParseUnverified(token.Raw, jwtv4.MapClaims{})
			if err != nil {
				return nil, err
			}
			claims, ok := t.Claims.(jwtv4.MapClaims)
			if !ok {
				return nil, fmt.Errorf("unexpected claims type")
			}
			// Check the "sub" claim
			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				// Handle the missing or invalid sub claim
				return nil, fmt.Errorf("invalid or missing 'sub' claim")
			}
			return jwks.Keyfunc(t)
		},
	})
}

func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}
