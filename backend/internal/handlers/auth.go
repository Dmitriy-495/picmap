package handlers

import (
	"context"
	"net/http"

	"picmap/internal/auth"
	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{DB: db}
}

// Login — POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, username, email, password_hash, is_superuser, is_admin, is_guide, is_owner, is_blocked, created_at
		 FROM users WHERE username = $1`,
		req.Username,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsSuperuser, &user.IsAdmin, &user.IsGuide, &user.IsOwner, &user.IsBlocked, &user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if user.IsBlocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "account blocked"})
		return
	}

	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateTokenWithOwner(user.ID, user.IsSuperuser, user.IsAdmin, user.IsGuide, user.IsOwner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// Me — GET /api/auth/me (требует JWT)
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var user models.User
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, username, email, is_superuser, is_admin, is_guide, is_owner, is_blocked, created_at
		 FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.IsSuperuser, &user.IsAdmin, &user.IsGuide, &user.IsOwner, &user.IsBlocked, &user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
