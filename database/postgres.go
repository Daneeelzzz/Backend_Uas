package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Gagal inisialisasi Postgres:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("❌ Gagal koneksi Postgres:", err)
	}

	log.Println("✅ Postgres Connected")
	return db
}