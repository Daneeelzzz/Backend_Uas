package main

import (
	"log"
	"tugas_uas/route"
	"tugas_uas/config"
	"tugas_uas/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	// Swagger Docs (Folder ini akan otomatis dibuat oleh perintah 'swag init')
	_ "tugas_uas/docs"
	"github.com/gofiber/swagger"
)

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description API untuk mengelola data prestasi, user, dan verifikasi dosen wali.
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Pengembang
// @contact.email admin@kampus.ac.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 1. Init Database
	dbContainer := config.InitDB()

	// 2. Seeder (Optional, nyalakan saat pertama kali saja)
	database.Seed(dbContainer.Postgres)

	// 3. Setup Fiber
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	// 4. Swagger Endpoint
	// Akses di: http://localhost:3000/swagger/index.html
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 5. Routes
	route.SetupRoutes(app, dbContainer)

	log.Println("🚀 Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}