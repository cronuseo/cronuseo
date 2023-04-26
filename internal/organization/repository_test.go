package organization

import (
	"context"
	"testing"
	"time"

	"github.com/shashimalcse/cronuseo/internal/test"
)

func TestRepository(t *testing.T) {

	db := test.DB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.MongoClient.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
}
