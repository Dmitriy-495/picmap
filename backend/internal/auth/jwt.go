package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — полезная нагрузка JWT
type Claims struct {
	UserID      int64 `json:"user_id"`
	IsSuperuser bool  `json:"is_superuser"`
	IsAdmin     bool  `json:"is_admin"`
	IsGuide     bool  `json:"is_guide"`
	jwt.RegisteredClaims
}

// GenerateToken создаёт JWT для пользователя
func GenerateToken(userID int64, isSuperuser, isAdmin, isGuide bool) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := Claims{
		UserID:      userID,
		IsSuperuser: isSuperuser,
		IsAdmin:     isAdmin,
		IsGuide:     isGuide,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 дней
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "picmap",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken проверяет JWT и возвращает Claims
func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
