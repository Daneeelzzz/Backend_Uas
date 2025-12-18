package service

import (
	"tugas_uas/app/repository"
	"tugas_uas/helper"

	"github.com/gofiber/fiber/v2"
)

type LecturerService interface {
	GetAllLecturers(c *fiber.Ctx) error
	GetMyAdvisees(c *fiber.Ctx) error // Dosen lihat mahasiswanya sendiri
}

type lecturerService struct {
	lecturerRepo repository.LecturerRepository
	studentRepo  repository.StudentRepository // Butuh akses ke student repo untuk cari anak wali
}

func NewLecturerService(l repository.LecturerRepository, s repository.StudentRepository) LecturerService {
	return &lecturerService{lecturerRepo: l, studentRepo: s}
}

// Admin: Lihat semua dosen
func (s *lecturerService) GetAllLecturers(c *fiber.Ctx) error {
	data, err := s.lecturerRepo.FindAll(c.Context())
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: data})
}

// Dosen: Lihat mahasiswa bimbingan SAYA (One-to-Many logic)
func (s *lecturerService) GetMyAdvisees(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	ctx := c.Context()

	// 1. Cari ID Lecturer (Table lecturers) dari UserID (Table users)
	lecturerID, err := s.lecturerRepo.GetLecturerIDByUserID(ctx, userID)
	if err != nil {
		return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "Profil Dosen tidak ditemukan"})
	}

	// 2. Cari semua mahasiswa yang punya advisor_id = lecturerID
	students, err := s.studentRepo.FindByAdvisorID(ctx, lecturerID)
	if err != nil {
		return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(helper.SuccessResponse{Success: true, Data: students})
}