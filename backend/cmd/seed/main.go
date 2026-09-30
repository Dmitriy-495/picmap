package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"picmap/internal/auth"
	"picmap/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using env vars")
	}

	ctx := context.Background()

	// Подключаемся к БД
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("❌ DB connect failed:", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to PostgreSQL")

	// Данные superuser
	username := "tda495"
	email := "tda495@picmap.ru"
	password := "23452345"

	// Хешируем пароль
	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal("❌ Hash password failed:", err)
	}

	// Проверяем, существует ли пользователь
	var exists bool
	err = pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		username,
	).Scan(&exists)
	if err != nil {
		log.Fatal("❌ Query failed:", err)
	}

	if exists {
		// Обновляем пароль
		_, err = pool.Exec(ctx,
			`UPDATE users 
			 SET password_hash = $1, is_superuser = TRUE, is_admin = TRUE, is_guide = TRUE 
			 WHERE username = $2`,
			hash, username,
		)
		if err != nil {
			log.Fatal("❌ Update failed:", err)
		}
		fmt.Printf("✅ Superuser '%s' updated\n", username)
	} else {
		// Создаём нового
		_, err = pool.Exec(ctx,
			`INSERT INTO users (username, email, password_hash, is_superuser, is_admin, is_guide) 
			 VALUES ($1, $2, $3, TRUE, TRUE, TRUE)`,
			username, email, hash,
		)
		if err != nil {
			log.Fatal("❌ Insert failed:", err)
		}
		fmt.Printf("✅ Superuser '%s' created\n", username)
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════")
	fmt.Println("✅ Seed complete!")
	fmt.Println("   Username: tda495")
	fmt.Println("   Password: 23452345")
	fmt.Println("═══════════════════════════════════")

	_ = os.Stdout
}
