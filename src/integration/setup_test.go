package integration

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/routes"
	"github.com/shashimalcse/Cronuseo/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"testing"
	"time"
)

func connectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME_TEST")
	db_link := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUsername, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(postgres.Open(db_link), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	migrateDB(db)
	config.DB = db
}

func migrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Group{}, &models.Organization{}, &models.Project{}, &models.Resource{},
		&models.ResourceAction{}, &models.ResourceRole{}, &models.GroupUser{}, &models.GroupUser{},
		&models.ResourceRoleToGroup{}, &models.ResourceRoleToUser{}, &models.ResourceRoleToResourceAction{},
		&models.ResourceRoleToResourceActionKey{})
	if err != nil {
		return
	}
}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("./../../.env"))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	connectDB()
	config.InitLogger()
	exitVal := m.Run()
	os.Exit(exitVal)

}

func StartServer(e *echo.Echo) {
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	routes.OrganizationRoutes(e)
	routes.ProjectRoutes(e)
	routes.ResourceRoutes(e)
	routes.ResourceActionRoutes(e)
	routes.ResourceRoutes(e)
	routes.GroupRoutes(e)
	routes.UserRoutes(e)
	routes.CheckRoutes(e)
}

func StopServer(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
