package service

import (
	"tugas_uas/app/repository"
	"tugas_uas/helper"

	"github.com/gofiber/fiber/v2"
)

type StudentService interface {
	GetAllStudents(c *fiber.Ctx) error
	AssignAdvisor(c *fiber.Ctx) error
}

type studentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

// Admin: Lihat semua mahasiswa + nama dosen walinya
func (s *studentService) GetAllStudents(c *fiber.Ctx) error {
	data, err := s.repo.FindAll(c.Context())
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: data})
}

// Admin: Assign Dosen Wali (Payload: {"advisor_id": "uuid-dosen"})
func (s *studentService) AssignAdvisor(c *fiber.Ctx) error {
	studentUserID := c.Params("id")
	
	var req struct {
		AdvisorID string `json:"advisor_id"`
	}
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Bad Input"}) }

	err := s.repo.UpdateAdvisor(c.Context(), studentUserID, req.AdvisorID)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Advisor Assigned"})
}