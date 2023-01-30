package cache

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
)

type PermissionCache interface {
	Set(context context.Context, key entity.Tuple, value string) error
	SetAPIKey(context context.Context, key string, value string) error
	Get(context context.Context, key entity.Tuple) (string, error)
	GetAPIKey(context context.Context, key string) (string, error)
	FlushAll(context context.Context) error
	DeleteAPIKey(context context.Context) error
}
