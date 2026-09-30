package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"picmap/internal/db"
	"picmap/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using environment variables")
	}

	// Режим Gin
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Подключаемся к БД
	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("❌ DB connect failed:", err)
	}
	defer pool.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Роутер
	r := gin.Default()

	// Health
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "picmap-backend",
			"version": "0.2.0",
		})
	})

	// Auth handlers
	authHandler := handlers.NewAuthHandler(pool)

	api := r.Group("/api")
	{
		// Публичные
		api.POST("/auth/login", authHandler.Login)

		// Защищённые
		protected := api.Group("")
		protected.Use(handlers.JWTMiddleware())
		{
			protected.GET("/auth/me", authHandler.Me)
		}
	}

	// Порт
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
