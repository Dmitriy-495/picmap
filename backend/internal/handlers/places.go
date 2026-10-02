package handlers

import (
	"context"
	"net/http"
	"strconv"

	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PlacesHandler — обработчики точек маршрута
type PlacesHandler struct {
	DB *pgxpool.Pool
}

func NewPlacesHandler(db *pgxpool.Pool) *PlacesHandler {
	return &PlacesHandler{DB: db}
}

// List — GET /api/routes/:id/places (публичный)
func (h *PlacesHandler) List(c *gin.Context) {
	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, route_id, name, lat, lng, COALESCE(description, ''), COALESCE(emoji, ''), order_index, created_at
		 FROM places WHERE route_id = $1 ORDER BY order_index`, routeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	places := []models.Place{}
	for rows.Next() {
		var p models.Place
		rows.Scan(&p.ID, &p.RouteID, &p.Name, &p.Lat, &p.Lng,
			&p.Description, &p.Emoji, &p.OrderIndex, &p.CreatedAt)
		places = append(places, p)
	}

	c.JSON(http.StatusOK, places)
}

// Create — POST /api/routes/:id/places (JWT + owner)
func (h *PlacesHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isSuperuser := c.GetBool("is_superuser")
	isAdmin := c.GetBool("is_admin")

	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	// Проверяем владельца
	var ownerID int64
	err = h.DB.QueryRow(context.Background(),
		`SELECT owner_id FROM routes WHERE id = $1`, routeID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}

	if ownerID != userID && !isSuperuser && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your route"})
		return
	}

	var req models.CreatePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var p models.Place
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO places (route_id, name, lat, lng, description, emoji, order_index)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, route_id, name, lat, lng, COALESCE(description, ''), COALESCE(emoji, ''), order_index, created_at`,
		routeID, req.Name, req.Lat, req.Lng, req.Description, req.Emoji, req.OrderIndex,
	).Scan(&p.ID, &p.RouteID, &p.Name, &p.Lat, &p.Lng,
		&p.Description, &p.Emoji, &p.OrderIndex, &p.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// Delete — DELETE /api/places/:id (JWT + owner)
func (h *PlacesHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isSuperuser := c.GetBool("is_superuser")
	isAdmin := c.GetBool("is_admin")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Проверяем владельца через JOIN
	var ownerID int64
	err = h.DB.QueryRow(context.Background(),
		`SELECT r.owner_id FROM places p JOIN routes r ON r.id = p.route_id WHERE p.id = $1`, id).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "place not found"})
		return
	}

	if ownerID != userID && !isSuperuser && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your place"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`DELETE FROM places WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "place deleted"})
}
