package models

import "time"

// Photo — фото, привязанное к маршруту или точке
type Photo struct {
	ID        int64     `json:"id"`
	RouteID   *int64    `json:"route_id,omitempty"`
	PlaceID   *int64    `json:"place_id,omitempty"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	FilePath  string    `json:"file_path"`
	URL       string    `json:"url"` // /uploads/...
	Title     string    `json:"title,omitempty"`
	Caption   string    `json:"caption,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdatePhotoRequest — обновление метаданных
type UpdatePhotoRequest struct {
	Title   *string `json:"title"`
	Caption *string `json:"caption"`
}
