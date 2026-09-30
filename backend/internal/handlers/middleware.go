package handlers

import (
	"net/http"
	"strings"

	"picmap/internal/auth"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware проверяет JWT в заголовке Authorization
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Ожидаем "Bearer <token>"
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Кладём данные в контекст
		c.Set("user_id", claims.UserID)
		c.Set("is_superuser", claims.IsSuperuser)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("is_guide", claims.IsGuide)

		c.Next()
	}
}
