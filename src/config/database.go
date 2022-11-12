package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db_link := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUsername, dbPass, dbHost, dbPort, dbName)
	db, err := sqlx.Connect("postgres", db_link)

	if err != nil {
		panic(err)
	}
	// migrateDB(db)
	DB = db
}

// func migrateDB(db *gorm.DB) {
// 	err := db.AutoMigrate(&models.User{}, &models.Group{}, &models.Organization{}, &models.Project{},
//	&models.Resource{},
// 		&models.ResourceAction{}, &models.ResourceRole{}, &models.GroupUser{}, &models.GroupUser{},
// 		&models.ResourceRoleToGroup{}, &models.ResourceRoleToUser{}, &models.ResourceRoleToResourceAction{},
// 		&models.ResourceRoleToResourceActionKey{})
// 	if err != nil {
// 		return
// 	}
// }
