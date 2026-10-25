package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"picmap/internal/db"

	"github.com/joho/godotenv"
)

type Superuser struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	IsOwner      bool   `json:"is_owner"`
	PasswordHash string `json:"password_hash"`
}

type Config struct {
	Superusers []Superuser `json:"superusers"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using env vars")
	}

	// Читаем config.json
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("❌ config.json not found:", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal("❌ Invalid config.json:", err)
	}

	if len(config.Superusers) == 0 {
		log.Fatal("❌ No superusers in config.json")
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("❌ DB connect failed:", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to PostgreSQL")

	for _, su := range config.Superusers {
		var exists bool
		err = pool.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
			su.Username,
		).Scan(&exists)
		if err != nil {
			log.Fatalf("❌ Query failed: %v", err)
		}

		if exists {
			_, err = pool.Exec(ctx,
				`UPDATE users 
				 SET email = $1, password_hash = $2,
				     is_superuser = TRUE, is_admin = TRUE, is_guide = TRUE,
				     is_owner = $3
				 WHERE username = $4`,
				su.Email, su.PasswordHash, su.IsOwner, su.Username,
			)
			if err != nil {
				log.Fatalf("❌ Update failed for %s: %v", su.Username, err)
			}
			fmt.Printf("✅ Superuser '%s' updated (owner=%v)\n", su.Username, su.IsOwner)
		} else {
			_, err = pool.Exec(ctx,
				`INSERT INTO users (username, email, password_hash, is_superuser, is_admin, is_guide, is_owner) 
				 VALUES ($1, $2, $3, TRUE, TRUE, TRUE, $4)`,
				su.Username, su.Email, su.PasswordHash, su.IsOwner,
			)
			if err != nil {
				log.Fatalf("❌ Insert failed for %s: %v", su.Username, err)
			}
			fmt.Printf("✅ Superuser '%s' created (owner=%v)\n", su.Username, su.IsOwner)
		}
	}

	fmt.Println("")
	fmt.Println("═══════════════════════════════════")
	fmt.Println("✅ Seed complete!")
	for _, su := range config.Superusers {
		ownerMark := ""
		if su.IsOwner {
			ownerMark = " 👑"
		}
		fmt.Printf("   %s%s\n", su.Username, ownerMark)
	}
	fmt.Println("═══════════════════════════════════")
}
