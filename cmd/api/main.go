package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"modmatch-be/internal/config"
	"modmatch-be/internal/database"
	"modmatch-be/internal/middleware"
	"modmatch-be/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.New(ctx, cfg.DB.DSN())
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	log.Info("database connected")

	// Wire dependencies here: pool -> repository -> service -> handler
	// e.g. tokenManager := auth.NewTokenManager(cfg.JWTSecret)

	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiV1 := router.Group("/api/v1")
	_ = apiV1 // register domain routes here, e.g. userHandler.RegisterRoutes(apiV1)

	addr := ":" + cfg.ServerPort
	log.Info("server started", "addr", addr)
	if err := router.Run(addr); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
