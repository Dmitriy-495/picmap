package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"picmap/internal/auth"
	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegisterHandler struct {
	DB *pgxpool.Pool
}

func NewRegisterHandler(db *pgxpool.Pool) *RegisterHandler {
	return &RegisterHandler{DB: db}
}

type SuperuserConfig struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	IsOwner  bool   `json:"is_owner"`
}

type Config struct {
	Superusers []SuperuserConfig `json:"superusers"`
}

// loadSuperusers — читает config.json
func loadSuperusers() (*Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// isSuperuser — проверяет, есть ли username в списке суперюзеров
func isSuperuser(username string) (bool, bool) {
	cfg, err := loadSuperusers()
	if err != nil {
		return false, false
	}
	for _, su := range cfg.Superusers {
		if su.Username == username {
			return true, su.IsOwner
		}
	}
	return false, false
}

// Register — POST /api/auth/register
// Регистрация по токену приглашения
func (h *RegisterHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 1. Проверяем приглашение
	var inv models.Invitation
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, token, created_by, used_by, expires_at, used_at, created_at
		 FROM invitations WHERE token = $1`, req.Token,
	).Scan(&inv.ID, &inv.Token, &inv.CreatedBy, &inv.UsedBy,
		&inv.ExpiresAt, &inv.UsedAt, &inv.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
		return
	}

	if inv.UsedAt != nil {
		c.JSON(http.StatusGone, gin.H{"error": "invitation already used"})
		return
	}

	if time.Now().After(inv.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "invitation expired"})
		return
	}

	// 2. Проверяем, что username свободен
	var exists bool
	err = h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		req.Username,
	).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	// 3. Хешируем пароль
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
		return
	}

	// 4. Проверяем, суперюзер ли это
	isSuper, isOwner := isSuperuser(req.Username)

	// 5. Создаём пользователя
	var user models.User
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password_hash, is_superuser, is_owner)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, username, email, is_superuser, is_admin, is_guide, is_owner, created_at`,
		req.Username, req.Username+"@picmap.ru", hash, isSuper, isOwner,
	).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.IsSuperuser, &user.IsAdmin, &user.IsGuide, &user.IsOwner, &user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 6. Помечаем приглашение использованным
	_, err = h.DB.Exec(context.Background(),
		`UPDATE invitations SET used_by = $1, used_at = NOW() WHERE id = $2`,
		user.ID, inv.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark invitation"})
		return
	}

	// 7. Генерируем JWT
	token, err := auth.GenerateToken(user.ID, user.IsSuperuser, user.IsAdmin, user.IsGuide)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		Token: token,
		User:  user,
	})
}
