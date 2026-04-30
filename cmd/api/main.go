package main

import (
	"log"

	"sluggo/config"
	docs "sluggo/docs"
	httphandler "sluggo/internal/handler/http"
	infracache "sluggo/internal/infrastructure/cache"
	"sluggo/internal/infrastructure/database"
	cacherepo "sluggo/internal/repository/cache"
	postgresrepo "sluggo/internal/repository/postgres"
	"sluggo/internal/usecase"
	"sluggo/pkg/shortener"

	"github.com/gin-gonic/gin"
)

// @title 		SlugGo API
// @version 	1.0
// @description A simple URL shortener service built with Go, Gin, and PostgreSQL.
// @BasePath 	/
// @schemes 	http, https
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}

	enc, err := shortener.New(cfg.HashSecretKey, cfg.SnowflakeNode)
	if err != nil {
		log.Fatalf("creating shortener: %v", err)
	}

	repo := postgresrepo.New(db)
	redisClient := infracache.Connect(cfg)
	cache := cacherepo.New(redisClient)
	uc := usecase.New(repo, cache, enc, cfg.AppBaseURL, cfg.CacheTTL)
	urlHandler := httphandler.NewURLHandler(uc)

	r := gin.Default()
	httphandler.RegisterRoutes(r, urlHandler)

	log.Printf("swagger docs available at %s/swagger/index.html", cfg.AppBaseURL)
	log.Printf("server listening on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
