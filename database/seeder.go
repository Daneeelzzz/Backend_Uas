package database

import (
	"database/sql"
	"log"
	"tugas_uas/utils"
)

func Seed(db *sql.DB) {
	log.Println("🌱 Memulai Smart Seeder...")

	// 1. Seed Roles
	roles := map[string]string{
		"Admin":      "Administrator Sistem",
		"Mahasiswa":  "Pengguna Mahasiswa",
		"Dosen Wali": "Dosen Pembimbing Akademik",
	}

	for name, desc := range roles {
		_, err := db.Exec("INSERT INTO roles (name, description) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING", name, desc)
		if err != nil {
			log.Printf("⚠️ Gagal seed role %s: %v", name, err)
		}
	}

	// 2. Ambil ID Role
	var adminRID, dosenRID, mhsRID string
	db.QueryRow("SELECT id FROM roles WHERE name = 'Admin'").Scan(&adminRID)
	db.QueryRow("SELECT id FROM roles WHERE name = 'Dosen Wali'").Scan(&dosenRID)
	db.QueryRow("SELECT id FROM roles WHERE name = 'Mahasiswa'").Scan(&mhsRID)

	pass, _ := utils.HashPassword("password123")

	// 3. Buat ADMIN
	var adminID string
	err := db.QueryRow(`
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) 
		VALUES ($1, $2, $3, $4, $5, true) 
		ON CONFLICT (username) DO UPDATE SET email=EXCLUDED.email 
		RETURNING id`, "admin_pusat", "admin@kampus.ac.id", pass, "Super Admin", adminRID).Scan(&adminID)
	if err == nil {
		log.Println("✅ Admin Created: admin_pusat")
	}

	// 4. Buat DOSEN WALI
	var dosenUserID, lecturerID string
	err = db.QueryRow(`
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) 
		VALUES ($1, $2, $3, $4, $5, true) 
		ON CONFLICT (username) DO UPDATE SET email=EXCLUDED.email 
		RETURNING id`, "dosen_budi", "budi@dosen.ac.id", pass, "Dr. Budi Santoso", dosenRID).Scan(&dosenUserID)
	
	if err == nil {
		// Insert ke tabel lecturers
		db.QueryRow(`
			INSERT INTO lecturers (user_id, lecturer_id, department) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (user_id) DO UPDATE SET lecturer_id=EXCLUDED.lecturer_id
			RETURNING id`, dosenUserID, "NIP19800101", "Teknik Informatika").Scan(&lecturerID)
		log.Println("✅ Dosen Created: dosen_budi")
	}

	// 5. Buat MAHASISWA (Langsung direlasikan ke Dosen Budi)
	var mhsUserID string
	err = db.QueryRow(`
		INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) 
		VALUES ($1, $2, $3, $4, $5, true) 
		ON CONFLICT (username) DO UPDATE SET email=EXCLUDED.email 
		RETURNING id`, "mhs_ani", "ani@student.ac.id", pass, "Ani Mahasiswi", mhsRID).Scan(&mhsUserID)
	
	if err == nil {
		// Insert ke tabel students DENGAN ADVISOR_ID
		_, err := db.Exec(`
			INSERT INTO students (user_id, student_id, program_study, advisor_id) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO NOTHING`, 
			mhsUserID, "NIM2025001", "S1 Informatika", lecturerID) // <--- Relasi disini
		
		if err != nil {
			log.Printf("❌ Gagal link mahasiswa ke dosen: %v", err)
		} else {
			log.Println("✅ Mahasiswa Created & Linked to Dosen Budi")
		}
	}

	log.Println("🌱 Seeding Selesai!")
}