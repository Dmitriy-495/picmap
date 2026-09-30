package models

import "time"

// Place — точка маршрута
type Place struct {
	ID          int64     `json:"id"`
	RouteID     int64     `json:"route_id"`
	Name        string    `json:"name"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Description string    `json:"description"`
	Emoji       string    `json:"emoji"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreatePlaceRequest — запрос на создание точки
type CreatePlaceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Lat         float64 `json:"lat" binding:"required"`
	Lng         float64 `json:"lng" binding:"required"`
	Description string  `json:"description"`
	Emoji       string  `json:"emoji"`
	OrderIndex  int     `json:"order_index"`
}
