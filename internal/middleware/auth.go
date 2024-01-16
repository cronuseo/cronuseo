package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shashimalcse/cronuseo/internal/check"
	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"go.uber.org/zap"
)

type MethodPath struct {
	Method   string
	Path     string
	Resource string
}

func Auth(cfg *config.Config, logger *zap.Logger, requiredPermissions map[MethodPath][]string, checkService check.Service) echo.MiddlewareFunc {

	jwks, err := keyfunc.Get(cfg.Auth.JWKS, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			logger.Error("There was an error with the jwt.KeyFunc", zap.Error(err))
		},
	})

	if err != nil {
		logger.Error("Failed to create JWKs from resource at the given URL.", zap.Error(err))
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				t, _, err := new(jwtv4.Parser).ParseUnverified(token.Raw, jwtv4.MapClaims{})
				if err != nil {
					return nil, err
				}

				claims, ok := t.Claims.(jwtv4.MapClaims)
				if !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected claims type")
				}

				//Extract and validate scopes
				sub, ok := claims["sub"].(string)
				if !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing sub claim")
				}

				methodPath := MethodPath{
					Method: c.Request().Method,
					Path:   c.Request().URL.Path,
				}
				logger.Debug("path", zap.String("path", methodPath.Path))
				pathMatched, err := regexp.MatchString("/api/v1/o/[^/]+/users/sync", methodPath.Path)
				if err != nil {
					return nil, err
				}
				if pathMatched {
					orgIdentifier := getOrgIdentifier(methodPath.Path)
					apiKey := c.Request().Header.Get("API_KEY")
					validated, _ := checkService.ValidateAPIKey(nil, orgIdentifier, apiKey)
					if !validated {
						return nil, echo.NewHTTPError(http.StatusUnauthorized, "insufficient permissions to invoke this endpoint")
					}
				} else {
					endpointPermissions, err := getPermissionsForMethodPath(methodPath, requiredPermissions)
					if err != nil {
						return nil, err
					}

					if !checkPermissions(sub, endpointPermissions, cfg, checkService) {
						logger.Debug("error while validating permissions")
						return nil, echo.NewHTTPError(http.StatusUnauthorized, "insufficient permissions to invoke this endpoint")
					}
				}
				key, keyErr := jwks.Keyfunc(t)
				if keyErr != nil {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "JWT key function error: %w", keyErr)
				}

				return key, nil
			}

			jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
				KeyFunc: keyFunc,
				ErrorHandlerWithContext: func(err error, c echo.Context) error {
					logger.Debug("error while validating token", zap.Error(err))
					if httpErr, ok := err.(*echo.HTTPError); ok {
						if internalErr, ok := httpErr.Internal.(*jwt.ValidationError); ok {
							return internalErr.Inner
						} else {
							return httpErr
						}
					}
					if validationErr, ok := err.(*jwt.ValidationError); ok {
						return echo.NewHTTPError(http.StatusUnauthorized, validationErr.Inner.Error())
					} else {
						return validationErr
					}
				},
			})
			return jwtMiddleware(next)(c)
		}
	}
}

// getScopesForMethodPath returns the scopes for a given method and path based on wildcard patterns.
func getPermissionsForMethodPath(methodPath MethodPath, requiredPermissions map[MethodPath][]string) ([]mongo_entity.Permission, error) {
	for pattern, permissions := range requiredPermissions {
		pathMatched, err := regexp.MatchString(pattern.Path, methodPath.Path)
		if err != nil {
			return nil, err
		}
		methodMatched := strings.EqualFold(pattern.Method, methodPath.Method) || pattern.Method == "*"

		if pathMatched && methodMatched {
			requiredPermissions := []mongo_entity.Permission{}
			for _, permission := range permissions {
				requiredPermissions = append(requiredPermissions, mongo_entity.Permission{Action: permission, Resource: pattern.Resource})
			}
			return requiredPermissions, nil
		}
	}
	return nil, fmt.Errorf("no matching scopes found for method and path: %s %s", methodPath.Method, methodPath.Path)
}

// checkScopes validates the required scopes are present for a given endpoint.
func checkPermissions(sub string, requiredPermissions []mongo_entity.Permission, cfg *config.Config, checkService check.Service) bool {

	for _, permission := range requiredPermissions {
		checkReq := check.CheckRequest{
			Identifier: sub,
			Action:     permission.Action,
			Resource:   permission.Resource,
		}
		allow, _ := checkService.Check(nil, cfg.RootOrganization.Name, checkReq, "nil", true)
		if !allow.Allowed {
			return false
		}
	}
	return true
}

func getOrgIdentifier(path string) string {

	re := regexp.MustCompile(`/api/v1/o/([^/]+)/users/sync`)
	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		org := matches[1]
		return org
	}
	return ""
}

func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}
