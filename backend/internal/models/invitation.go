package models

import "time"

// Invitation — приглашение по ссылке
type Invitation struct {
	ID        int64      `json:"id"`
	Token     string     `json:"token"`
	CreatedBy int64      `json:"created_by"`
	UsedBy    *int64     `json:"used_by,omitempty"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// CreateInvitationRequest — запрос на создание приглашения
type CreateInvitationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// CreateInvitationResponse — ответ с готовой ссылкой
type CreateInvitationResponse struct {
	Invitation Invitation `json:"invitation"`
	URL        string     `json:"url"`
}

// RegisterRequest — регистрация по токену
type RegisterRequest struct {
	Token    string `json:"token" binding:"required"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse — ответ после регистрации
type RegisterResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
