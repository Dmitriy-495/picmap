package models

import "time"

// Comment — комментарий в ленте маршрута
type Comment struct {
	ID        int64     `json:"id"`
	RouteID   int64     `json:"route_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`  // JOIN с users
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CreateCommentRequest — запрос на создание
type CreateCommentRequest struct {
	Text string `json:"text" binding:"required"`
}

// UpdateCommentRequest — запрос на обновление
type UpdateCommentRequest struct {
	Text string `json:"text" binding:"required"`
}
