package main

import (
	"context"
	"flag"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/open-policy-agent/opa/rego"
	_ "github.com/shashimalcse/cronuseo/docs"
	"github.com/shashimalcse/cronuseo/internal/check"
	"github.com/shashimalcse/cronuseo/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Version = "1.0.0"

// Default config flag.
var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

// @title          Cronuseo API
// @version        1.0
// @description    This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8081
// @BasePath /api/v1
func main() {

	logger := InitializeLogger()

	flag.Parse()

	// Load configurations for db, keto and redis.
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal("Error while loading config.", zap.String("config_file", flag.Lookup("config").Value.String()))
		os.Exit(-1)
	}

	// Mongo client.
	credential := options.Credential{
		Username: cfg.MongoUser,
		Password: cfg.MongoPassword,
	}
	mongo_client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Mongo).SetAuth(credential))
	if err != nil {
		panic(err)
	}

	mongo_db := mongo_client.Database(cfg.MongoDBName)

	r := rego.New(
		rego.Query("x = data.example.allow"),
		rego.Load([]string{cfg.RBACPolicy}, nil))
	ctx := context.Background()
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		panic(err)
	}
	e := buildHandler(cfg, logger, mongo_db, query)
	logger.Info("Starting server", zap.String("server_endpoint", cfg.Check_API))
	e.Logger.Fatal(e.Start(cfg.Check_API))

}

// buildHandler builds the echo router.
func buildHandler(
	cfg *config.Config, // Config
	logger *zap.Logger, // Logger
	mongodb *mongo.Database, // Mongo monitoring client
	query rego.PreparedEvalQuery,
) *echo.Echo {

	router := echo.New()

	//  Set up CORS.
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "API_KEY"}, // API_KEY is used for permission checking SDKs
		AllowOrigins:     []string{"http://localhost:3000"},
	}))
	// Echo logger middleware.
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}; method=${method}; uri=${uri}; status=${status};\n",
	}))

	// API endpoints.
	rg := router.Group("/api/v1")

	// Here we register all the handlers. Each handler handle jwt validation separately.
	check.RegisterHandlers(rg, check.NewService(check.NewRepository(mongodb), logger, query))

	return router
}

func InitializeLogger() *zap.Logger {

	zap_config := zap.NewProductionEncoderConfig()
	zap_config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zap_config)
	consoleEncoder := zapcore.NewConsoleEncoder(zap_config)
	logFile, _ := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
