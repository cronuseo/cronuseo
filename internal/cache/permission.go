package cache

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
)

type PermissionCache interface {
	Set(context context.Context, key entity.Tuple, value string) error
	Get(context context.Context, key entity.Tuple) (string, error)
	FlushAll(context context.Context) error
}
