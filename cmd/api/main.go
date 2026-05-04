package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"sluggo/config"
	docs "sluggo/docs"
	httphandler "sluggo/internal/handler/http"
	infracache "sluggo/internal/infrastructure/cache"
	"sluggo/internal/infrastructure/database"
	cacherepo "sluggo/internal/repository/cache"
	postgresrepo "sluggo/internal/repository/postgres"
	"sluggo/internal/usecase"
	"sluggo/pkg/logger"
	"sluggo/pkg/shortener"

	"github.com/gin-gonic/gin"
)

// @title 		SlugGo API
// @version 	1.0
// @description A simple URL shortener service built with Go, Gin, and PostgreSQL.
// @BasePath 	/
// @schemes 	http, https
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	// Set up logging based on environment
	logger.SetLogger(cfg.AppEnv)

	// Set Swagger info
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Initialize database connection
	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("connecting to database", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize URL shortener
	enc, err := shortener.New(cfg.HashSecretKey, cfg.SnowflakeNode)
	if err != nil {
		slog.Error("creating shortener", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize repositories, use case, and HTTP handler
	repo := postgresrepo.New(db)
	redisClient := infracache.Connect(cfg)
	cache := cacherepo.New(redisClient)
	uc := usecase.New(repo, cache, enc, cfg.AppBaseURL, cfg.CacheTTL)
	urlHandler := httphandler.NewURLHandler(uc)

	// Set up Gin router and register routes
	r := gin.Default()
	httphandler.RegisterRoutes(r, urlHandler)

	slog.Info("swagger docs available", slog.String("url", fmt.Sprintf("%s/swagger/index.html", cfg.AppBaseURL)))
	slog.Info("server listening", slog.String("port", cfg.AppPort))
	if err := r.Run(":" + cfg.AppPort); err != nil {
		slog.Error("starting server", slog.Any("error", err))
		os.Exit(1)
	}
}
