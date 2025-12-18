package route

import (
	"tugas_uas/app/repository"
	"tugas_uas/app/service"
	"tugas_uas/config"
	"tugas_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *config.DatabaseContainer) {
	// === 1. REPOSITORIES ===
	userRepo := repository.NewUserRepository(db.Postgres)
	lecturerRepo := repository.NewLecturerRepository(db.Postgres) // BARU
	studentRepo := repository.NewStudentRepository(db.Postgres)   // BARU
	achRepo := repository.NewAchievementRepository(db.Postgres, db.Mongo)

	// === 2. SERVICES ===
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	
	// Inject StudentRepo ke LecturerService (untuk fitur get advisees)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo) 
	studentService := service.NewStudentService(studentRepo)

	achService := service.NewAchievementService(achRepo, userRepo)

	// === 3. ROUTES ===
	api := app.Group("/api/v1")

	// Auth
	api.Post("/auth/login", authService.Login)

	// Protected
	app.Use(middleware.AuthRequired())

	// --- MODULE: USERS ---
	api.Get("/users/me", userService.GetMyProfile)
	api.Put("/users/me", userService.UpdateMyProfile)
	
	adm := api.Group("/admin", middleware.PermissionCheck("Admin"))
	adm.Post("/users", userService.CreateUser) // Admin Create User
	
	// --- MODULE: STUDENTS (SRS 5.5) ---
	// Admin melihat semua mahasiswa
	api.Get("/students", middleware.PermissionCheck("Admin"), studentService.GetAllStudents)
	// Admin set dosen wali
	api.Put("/students/:id/advisor", middleware.PermissionCheck("Admin"), studentService.AssignAdvisor)

	// --- MODULE: LECTURERS (SRS 5.5) ---
	// Admin melihat semua dosen
	api.Get("/lecturers", middleware.PermissionCheck("Admin"), lecturerService.GetAllLecturers)
	
	// Dosen Wali melihat mahasiswanya sendiri (RELASI ONE TO MANY)
	api.Get("/lecturer/my-students", middleware.PermissionCheck("Dosen Wali"), lecturerService.GetMyAdvisees)

	// --- MODULE: ACHIEVEMENTS ---
	ach := api.Group("/achievements")
	ach.Post("/", middleware.PermissionCheck("Mahasiswa"), achService.CreateAchievement)
	ach.Get("/me", middleware.PermissionCheck("Mahasiswa"), achService.GetMyHistory)
	ach.Put("/:id", middleware.PermissionCheck("Mahasiswa"), achService.UpdateAchievement)
	ach.Delete("/:id", middleware.PermissionCheck("Mahasiswa"), achService.DeleteAchievement)
	ach.Post("/:id/submit", middleware.PermissionCheck("Mahasiswa"), achService.SubmitForVerification)

	ach.Put("/:id/verify", middleware.PermissionCheck("Dosen Wali"), achService.VerifyAchievement)
	ach.Put("/:id/reject", middleware.PermissionCheck("Dosen Wali"), achService.RejectAchievement)

	admAch := api.Group("/admin/achievements", middleware.PermissionCheck("Admin"))
	admAch.Get("/", achService.GetAllAchievements)
	admAch.Put("/:id/finalize", achService.AdminFinalize)
	
	api.Get("/reports/stats", middleware.PermissionCheck("Admin"), achService.GetStatistics)
}