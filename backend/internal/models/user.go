package models

import "time"

// User — пользователь системы
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsSuperuser  bool      `json:"is_superuser"`
	IsAdmin      bool      `json:"is_admin"`
	IsGuide      bool      `json:"is_guide"`
	CreatedAt    time.Time `json:"created_at"`
}

// LoginRequest — запрос на логин
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse — ответ при успешном логине
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
