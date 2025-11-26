package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver pgx
)

func ConnectPostgres() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka driver: %w", err)
	}

	// Konfigurasi Pool agar performa cepat
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Tes Ping
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping postgres: %w", err)
	}

	log.Println("✅ Terhubung ke PostgreSQL (Raw SQL)")
	return db, nil
}