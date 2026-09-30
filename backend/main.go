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
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using env vars")
	}

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("❌ DB connect failed:", err)
	}
	defer pool.Close()
	log.Println("✅ Connected to PostgreSQL")

	r := gin.Default()

	// Health
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "picmap-backend",
			"version": "0.3.0",
		})
	})

	// Handlers
	authHandler := handlers.NewAuthHandler(pool)
	routesHandler := handlers.NewRoutesHandler(pool)

	api := r.Group("/api")
	{
		// Публичные
		api.POST("/auth/login", authHandler.Login)
		api.GET("/routes", routesHandler.List)   // список публичный

		// Защищённые (JWT)
		protected := api.Group("")
		protected.Use(handlers.JWTMiddleware())
		{
			protected.GET("/auth/me", authHandler.Me)
			protected.GET("/routes/:id", routesHandler.Get)
			protected.POST("/routes", routesHandler.Create)
			protected.DELETE("/routes/:id", routesHandler.Delete)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
