package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"picmap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PhotosHandler — обработчики фото
type PhotosHandler struct {
	DB        *pgxpool.Pool
	UploadDir string
}

func NewPhotosHandler(db *pgxpool.Pool) *PhotosHandler {
	dir := os.Getenv("UPLOAD_DIR")
	if dir == "" {
		dir = "./uploads"
	}
	return &PhotosHandler{DB: db, UploadDir: dir}
}

// Upload — POST /api/photos/upload (multipart, JWT)
// Параметры формы:
//   - file     (обязательно)
//   - route_id (опционально)
//   - place_id (опционально)
//   - title    (опционально)
//   - caption  (опционально)
func (h *PhotosHandler) Upload(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Файл
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	// Лимит 20 MB
	if file.Size > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 20MB)"})
		return
	}

	// Проверяем расширение
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
		return
	}

	// Привязка
	var routeID, placeID *int64
	if v := c.PostForm("route_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			routeID = &id
		}
	}
	if v := c.PostForm("place_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			placeID = &id
		}
	}

	if routeID == nil && placeID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "route_id or place_id required"})
		return
	}

	// Создаём папку для загрузок
	if err := os.MkdirAll(h.UploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
		return
	}

	// Генерируем имя файла
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), "photo", ext)
	filePath := filepath.Join(h.UploadDir, filename)

	// Сохраняем файл
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
		return
	}

	// URL для отдачи
	url := "/uploads/" + filename

	// Записываем в БД
	var p models.Photo
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO photos (route_id, place_id, user_id, file_path, title, caption)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, route_id, place_id, user_id, file_path, COALESCE(title, ''), COALESCE(caption, ''), created_at`,
		routeID, placeID, userID, filePath,
		c.PostForm("title"), c.PostForm("caption"),
	).Scan(&p.ID, &p.RouteID, &p.PlaceID, &p.UserID,
		&p.FilePath, &p.Title, &p.Caption, &p.CreatedAt)

	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.URL = url
	h.DB.QueryRow(context.Background(),
		`SELECT username FROM users WHERE id = $1`, userID).Scan(&p.Username)

	c.JSON(http.StatusCreated, p)
}

// ListByRoute — GET /api/routes/:id/photos (публичный)
func (h *PhotosHandler) ListByRoute(c *gin.Context) {
	routeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route id"})
		return
	}

	h.listPhotos(c, "route_id", routeID)
}

// ListByPlace — GET /api/places/:id/photos (публичный)
func (h *PhotosHandler) ListByPlace(c *gin.Context) {
	placeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place id"})
		return
	}

	h.listPhotos(c, "place_id", placeID)
}

// listPhotos — общая логика
func (h *PhotosHandler) listPhotos(c *gin.Context, field string, id int64) {
	query := fmt.Sprintf(
		`SELECT p.id, p.route_id, p.place_id, p.user_id, u.username,
		        p.file_path, COALESCE(p.title, ''), COALESCE(p.caption, ''), p.created_at
		 FROM photos p
		 JOIN users u ON u.id = p.user_id
		 WHERE p.%s = $1
		 ORDER BY p.created_at DESC`, field)

	rows, err := h.DB.Query(context.Background(), query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	photos := []models.Photo{}
	for rows.Next() {
		var p models.Photo
		rows.Scan(&p.ID, &p.RouteID, &p.PlaceID, &p.UserID, &p.Username,
			&p.FilePath, &p.Title, &p.Caption, &p.CreatedAt)
		p.URL = "/uploads/" + filepath.Base(p.FilePath)
		photos = append(photos, p)
	}

	c.JSON(http.StatusOK, photos)
}

// Delete — DELETE /api/photos/:id (JWT + owner)
func (h *PhotosHandler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	isSuperuser := c.GetBool("is_superuser")
	isAdmin := c.GetBool("is_admin")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var ownerID int64
	var filePath string
	err = h.DB.QueryRow(context.Background(),
		`SELECT user_id, file_path FROM photos WHERE id = $1`, id).Scan(&ownerID, &filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}

	if ownerID != userID && !isSuperuser && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your photo"})
		return
	}

	// Удаляем из БД
	_, err = h.DB.Exec(context.Background(),
		`DELETE FROM photos WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Удаляем файл
	os.Remove(filePath)

	c.JSON(http.StatusOK, gin.H{"message": "photo deleted"})
}
