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
	_ "tugas_uas/docs" // Import docs Swagger
	"github.com/gofiber/swagger"
)

// @title Sistem Prestasi API
// @version 1.0
// @description API Pelaporan Prestasi Mahasiswa
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load()
	db := config.InitDB()
	database.Seed(db.Postgres) // Smart Seeder

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault)
	route.SetupRoutes(app, db)

	log.Println("🚀 Server running on port 3000")
	log.Fatal(app.Listen(":3000"))
}