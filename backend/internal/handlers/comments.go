package handlers

import (
	"context"
	"net/http"
	"strconv"

	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommentsHandler — обработчики комментариев
type CommentsHandler struct {
	DB *pgxpool.Pool
}

func NewCommentsHandler(db *pgxpool.Pool) *CommentsHandler {
	return &CommentsHandler{DB: db}
}

// List — GET /api/routes/:id/comments (публичный)
func (h *CommentsHandler) List(c *gin.Context) {
	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT c.id, c.route_id, c.user_id, u.username, c.text, c.created_at, c.updated_at
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.route_id = $1
		 ORDER BY c.created_at ASC`, routeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var cm models.Comment
		rows.Scan(&cm.ID, &cm.RouteID, &cm.UserID, &cm.Username,
			&cm.Text, &cm.CreatedAt, &cm.UpdatedAt)
		comments = append(comments, cm)
	}

	c.JSON(http.StatusOK, comments)
}

// Create — POST /api/routes/:id/comments (JWT)
func (h *CommentsHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var cm models.Comment
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO comments (route_id, user_id, text)
		 VALUES ($1, $2, $3)
		 RETURNING id, route_id, user_id, text, created_at, updated_at`,
		routeID, userID, req.Text,
	).Scan(&cm.ID, &cm.RouteID, &cm.UserID, &cm.Text, &cm.CreatedAt, &cm.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Подтягиваем username
	h.DB.QueryRow(context.Background(),
		`SELECT username FROM users WHERE id = $1`, userID).Scan(&cm.Username)

	c.JSON(http.StatusCreated, cm)
}

// Update — PUT /api/comments/:id (JWT + owner)
func (h *CommentsHandler) Update(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isSuperuser := c.GetBool("is_superuser")
	isAdmin := c.GetBool("is_admin")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Проверяем владельца
	var ownerID int64
	err = h.DB.QueryRow(context.Background(),
		`SELECT user_id FROM comments WHERE id = $1`, id).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	if ownerID != userID && !isSuperuser && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your comment"})
		return
	}

	var cm models.Comment
	err = h.DB.QueryRow(context.Background(),
		`UPDATE comments SET text = $1, updated_at = NOW()
		 WHERE id = $2
		 RETURNING id, route_id, user_id, text, created_at, updated_at`,
		req.Text, id,
	).Scan(&cm.ID, &cm.RouteID, &cm.UserID, &cm.Text, &cm.CreatedAt, &cm.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.DB.QueryRow(context.Background(),
		`SELECT username FROM users WHERE id = $1`, cm.UserID).Scan(&cm.Username)

	c.JSON(http.StatusOK, cm)
}

// Delete — DELETE /api/comments/:id (JWT + owner)
func (h *CommentsHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isSuperuser := c.GetBool("is_superuser")
	isAdmin := c.GetBool("is_admin")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var ownerID int64
	err = h.DB.QueryRow(context.Background(),
		`SELECT user_id FROM comments WHERE id = $1`, id).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	if ownerID != userID && !isSuperuser && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your comment"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}
