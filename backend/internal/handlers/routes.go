package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoutesHandler — обработчики маршрутов
type RoutesHandler struct {
	DB *pgxpool.Pool
}

func NewRoutesHandler(db *pgxpool.Pool) *RoutesHandler {
	return &RoutesHandler{DB: db}
}

// List — GET /api/routes
func (h *RoutesHandler) List(c *gin.Context) {
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, owner_id, title, COALESCE(description, ''), start_date, end_date, cover_photo_id, created_at
		 FROM routes
		 ORDER BY 
		   CASE WHEN start_date > NOW() THEN 0 ELSE 1 END,
		   start_date DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var r models.Route
		if err := rows.Scan(
			&r.ID, &r.OwnerID, &r.Title, &r.Description,
			&r.StartDate, &r.EndDate, &r.CoverPhotoID, &r.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		r.Status = r.ComputeStatus()
		routes = append(routes, r)
	}

	c.JSON(http.StatusOK, routes)
}

// Get — GET /api/routes/:id
func (h *RoutesHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var r models.Route
	err = h.DB.QueryRow(context.Background(),
		`SELECT id, owner_id, title, COALESCE(description, ''), start_date, end_date, cover_photo_id, created_at
		 FROM routes WHERE id = $1`, id,
	).Scan(
		&r.ID, &r.OwnerID, &r.Title, &r.Description,
		&r.StartDate, &r.EndDate, &r.CoverPhotoID, &r.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	r.Status = r.ComputeStatus()

	// Подгружаем точки
	rows, err := h.DB.Query(context.Background(),
		`SELECT id, route_id, name, lat, lng, COALESCE(description, ''), COALESCE(emoji, ''), order_index, created_at
		 FROM places WHERE route_id = $1 ORDER BY order_index`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p models.Place
			rows.Scan(&p.ID, &p.RouteID, &p.Name, &p.Lat, &p.Lng,
				&p.Description, &p.Emoji, &p.OrderIndex, &p.CreatedAt)
			r.Places = append(r.Places, p)
		}
	}

	c.JSON(http.StatusOK, r)
}

// Create — POST /api/routes (JWT + guide)
func (h *RoutesHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req models.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Парсим даты
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (YYYY-MM-DD)"})
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (YYYY-MM-DD)"})
		return
	}

	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	var r models.Route
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO routes (owner_id, title, description, start_date, end_date)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, owner_id, title, COALESCE(description, ''), start_date, end_date, cover_photo_id, created_at`,
		userID, req.Title, req.Description, startDate, endDate,
	).Scan(
		&r.ID, &r.OwnerID, &r.Title, &r.Description,
		&r.StartDate, &r.EndDate, &r.CoverPhotoID, &r.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r.Status = r.ComputeStatus()

	c.JSON(http.StatusCreated, r)
}

// Delete — DELETE /api/routes/:id (JWT + owner/admin)
func (h *RoutesHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isSuperuser := c.GetBool("is_superuser")
	isAdmin := c.GetBool("is_admin")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Проверяем владельца
	var ownerID int64
	err = h.DB.QueryRow(context.Background(),
		`SELECT owner_id FROM routes WHERE id = $1`, id).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}

	if ownerID != userID && !isSuperuser && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your route"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`DELETE FROM routes WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "route deleted"})
}
