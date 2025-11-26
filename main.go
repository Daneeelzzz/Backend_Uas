package main

import (
	"log"
	"pbl-backend/database" // Sesuaikan dengan module name di go.mod kamu

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment sistem/default")
	}

	// 2. Koneksi Postgres (Raw SQL)
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("❌ Fatal: Gagal koneksi database: %v", err)
	}
	// Jangan lupa tutup koneksi saat aplikasi mati
	defer db.Close()

	// 3. Tes Query Sederhana (Bukti kita pakai Raw SQL)
	var now string
	// Kita menjalankan query SQL manual "SELECT NOW()"
	err = db.QueryRow("SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal("Gagal eksekusi query test:", err)
	}

	log.Printf("Sistem siap! Waktu Database saat ini: %s", now)
}