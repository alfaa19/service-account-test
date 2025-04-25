package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alfaa19/service-account-test/config"
	"github.com/alfaa19/service-account-test/internal/handler"
	"github.com/alfaa19/service-account-test/internal/repository"
	"github.com/alfaa19/service-account-test/internal/routes"
	"github.com/alfaa19/service-account-test/internal/service"
	logger "github.com/alfaa19/service-account-test/pkg/logrus"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	customLogger, err := logger.NewLogger(cfg.GetLoggerConfig())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Connect to database
	err = cfg.OpenDatabase()
	if err != nil {
		customLogger.LogOperation(context.Background(), "main", "error", map[string]interface{}{
			"error": "Failed to connect to database: " + err.Error(),
		})

	}
	defer cfg.DBConnection.Close()
	if err := cfg.DBConnection.Ping(); err != nil {
		customLogger.Fatal("Failed to ping database: ", err)
	}
	// Initialize dependencies
	repo := repository.NewRepository(cfg.DBConnection, customLogger)
	svc := service.NewService(repo, customLogger)
	h := handler.NewAccountHandler(svc, customLogger)

	// Initialize Echo
	e := echo.New()

	// Setup routes
	routes.NewRouter(h, e)

	// Start server
	go func() {
		if err := e.Start(cfg.GetServerAddress()); err != nil {
			customLogger.LogOperation(context.Background(), "main", "error", map[string]interface{}{
				"error": "Server failed to start: " + err.Error(),
			})
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		customLogger.LogOperation(ctx, "main", "error", map[string]interface{}{
			"error": "Server shutdown failed: " + err.Error(),
		})
	}

	customLogger.LogOperation(ctx, "main", "success", map[string]interface{}{
		"message": "Server shutdown completed",
	})
}
