package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/open-policy-agent/opa/rego"
	_ "github.com/shashimalcse/cronuseo/docs"
	"github.com/shashimalcse/cronuseo/internal/check"
	"github.com/shashimalcse/cronuseo/internal/config"
	db "github.com/shashimalcse/cronuseo/internal/db/mongo"
	"github.com/shashimalcse/cronuseo/internal/logger"
	"github.com/shashimalcse/cronuseo/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	flag.Parse()

	// Load configurations.
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatal("Error while loading config.", zap.String("config_file", flag.Lookup("config").Value.String()))
		os.Exit(-1)
	}

	// Set up logger.
	logger := logger.Init(cfg)

	// Mongo client.
	mongodb := db.Init(cfg, logger)

	// OPA policy.
	r := rego.New(
		rego.Query("x = data.example.allow"),
		rego.Load([]string{cfg.OPA.RBAC}, nil))
	ctx := context.Background()
	query, err := r.PrepareForEval(ctx)

	if err != nil {
		log.Fatal("Error while prepare rego policy.")
		os.Exit(-1)
	}

	// Start the REST server
	go func() {
		e := BuildHandler(cfg, logger, mongodb, query)
		logger.Info("Starting REST server", zap.String("REST server_endpoint", cfg.Endpoint.Check_REST))
		e.Logger.Fatal(e.Start(cfg.Endpoint.Check_REST))
	}()

	if cfg.Endpoint.Check_GRPC != "" {
		// Start the gRPC server
		go func() {
			lis, err := net.Listen("tcp", cfg.Endpoint.Check_GRPC)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			service := check.NewGrpcService(check.NewService(check.NewRepository(mongodb), logger, query), logger)
			s := grpc.NewServer()
			proto.RegisterCheckServer(s, service)
			logger.Info("Starting GRPC server", zap.String("GRPC server_endpoint", cfg.Endpoint.Check_GRPC))
			log.Fatal(s.Serve(lis))
		}()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info("Shutting down servers...")
	os.Exit(0)

}

// buildHandler builds the echo router.
func BuildHandler(
	cfg *config.Config, // Config
	logger *zap.Logger, // Logger
	mongodb *db.MongoDB, // MongoDB
	query rego.PreparedEvalQuery,
) *echo.Echo {

	router := echo.New()

	// Set up CORS.
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

	// Here we register all the handlers.
	check.RegisterHandlers(rg, check.NewService(check.NewRepository(mongodb), logger, query))

	return router
}
