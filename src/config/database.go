package config

import (
	"fmt"
	"os"

	"github.com/shashimalcse/Cronuseo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db_link := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUsername, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(postgres.Open(db_link), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.Organization{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Resource{})
	db.AutoMigrate(&models.ResourceAction{})
	db.AutoMigrate(&models.ResourceRole{})
	DB = db
}
