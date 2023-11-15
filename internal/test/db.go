package test

import (
	"log"
	"testing"

	"github.com/shashimalcse/cronuseo/internal/config"
	db "github.com/shashimalcse/cronuseo/internal/db/mongo"
)

func DB(t *testing.T) *db.MongoDB {

	cfg, err := config.Load("../../config/run-test.yml")
	if err != nil {
		log.Fatal("Error while loading config for test.")
	}
	logger := InitLogger()
	mongo, err := db.Init(cfg, logger)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB client")
	}
	return mongo

}
