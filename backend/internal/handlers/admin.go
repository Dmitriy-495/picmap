package handlers

import (
	"context"
	"net/http"

	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminHandler struct {
	DB *pgxpool.Pool
}

func NewAdminHandler(db *pgxpool.Pool) *AdminHandler {
	return &AdminHandler{DB: db}
}

// ListUsers — GET /api/admin/users
// Доступ: только owner
func (h *AdminHandler) ListUsers(c *gin.Context) {
	isOwner := c.GetBool("is_owner")

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can view users"})
		return
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, username, email, is_superuser, is_admin, is_guide, is_owner, is_blocked, created_at
		 FROM users
		 ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Email,
			&u.IsSuperuser, &u.IsAdmin, &u.IsGuide, &u.IsOwner, &u.IsBlocked, &u.CreatedAt)
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}
