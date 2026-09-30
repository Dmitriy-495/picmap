package models

import "time"

// Route — маршрут путешествия
type Route struct {
	ID           int64      `json:"id"`
	OwnerID      int64      `json:"owner_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      time.Time  `json:"end_date"`
	CoverPhotoID *int64     `json:"cover_photo_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	Status       string     `json:"status,omitempty"`   // вычисляемое: будет/сейчас/было
	Places       []Place    `json:"places,omitempty"`
}

// CreateRouteRequest — запрос на создание маршрута
type CreateRouteRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	StartDate   string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate     string `json:"end_date" binding:"required"`   // YYYY-MM-DD
}

// UpdateRouteRequest — запрос на обновление
type UpdateRouteRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

// ComputeStatus — вычисляет статус по датам
func (r *Route) ComputeStatus() string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	if r.StartDate.After(today) {
		return "будет"
	}
	if r.EndDate.Before(today) {
		return "было"
	}
	return "сейчас"
}
