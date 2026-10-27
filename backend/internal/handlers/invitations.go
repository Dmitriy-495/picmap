package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvitationsHandler struct {
	DB *pgxpool.Pool
}

func NewInvitationsHandler(db *pgxpool.Pool) *InvitationsHandler {
	return &InvitationsHandler{DB: db}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Create — POST /api/admin/invitations (owner only)
func (h *InvitationsHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isOwner := c.GetBool("is_owner")

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can create invitations"})
		return
	}

	var req models.CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	var inv models.Invitation
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO invitations (token, email, created_by, expires_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, token, email, created_by, used_by, expires_at, used_at, created_at`,
		token, req.Email, userID, expiresAt,
	).Scan(
		&inv.ID, &inv.Token, &inv.Email, &inv.CreatedBy, &inv.UsedBy,
		&inv.ExpiresAt, &inv.UsedAt, &inv.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	c.JSON(http.StatusCreated, models.CreateInvitationResponse{
		Invitation: inv,
		URL:        frontendURL + "/invite/" + inv.Token,
	})
}

// List — GET /api/admin/invitations (owner only)
func (h *InvitationsHandler) List(c *gin.Context) {
	isOwner := c.GetBool("is_owner")

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can view invitations"})
		return
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, token, COALESCE(email, ''), created_by, used_by, expires_at, used_at, created_at
		 FROM invitations
		 ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	invitations := []models.Invitation{}
	for rows.Next() {
		var inv models.Invitation
		rows.Scan(&inv.ID, &inv.Token, &inv.Email, &inv.CreatedBy, &inv.UsedBy,
			&inv.ExpiresAt, &inv.UsedAt, &inv.CreatedAt)
		invitations = append(invitations, inv)
	}

	c.JSON(http.StatusOK, invitations)
}

// Delete — DELETE /api/admin/invitations/:id (owner only)
func (h *InvitationsHandler) Delete(c *gin.Context) {
	isOwner := c.GetBool("is_owner")

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete invitations"})
		return
	}

	id := c.Param("id")

	_, err := h.DB.Exec(context.Background(),
		`DELETE FROM invitations WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation deleted"})
}

// GetByToken — GET /api/invitations/:token (public)
func (h *InvitationsHandler) GetByToken(c *gin.Context) {
	token := c.Param("token")

	var inv models.Invitation
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, token, COALESCE(email, ''), created_by, used_by, expires_at, used_at, created_at
		 FROM invitations WHERE token = $1`, token,
	).Scan(&inv.ID, &inv.Token, &inv.Email, &inv.CreatedBy, &inv.UsedBy,
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

	c.JSON(http.StatusOK, inv)
}
