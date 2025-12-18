package service

import (
	"context"
	"time"
	"tugas_uas/app/model"
	"tugas_uas/app/repository"
	"tugas_uas/helper"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error // NEW: Edit Draft
	DeleteAchievement(c *fiber.Ctx) error // Delete Draft
	GetMyHistory(c *fiber.Ctx) error
	
	SubmitForVerification(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	
	GetAllAchievements(c *fiber.Ctx) error
	GetStatistics(c *fiber.Ctx) error
	AdminFinalize(c *fiber.Ctx) error // NEW: Admin ACC
}

type achService struct {
	achRepo  repository.AchievementRepository
	userRepo repository.UserRepository
}

func NewAchievementService(a repository.AchievementRepository, u repository.UserRepository) AchievementService {
	return &achService{achRepo: a, userRepo: u}
}

// ... (Method Create Sama, Skip) ...
func (s *achService) CreateAchievement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req model.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Bad Input"}) }

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	mongoData := &model.AchievementMongo{
		ID:              primitive.NewObjectID(),
		StudentRefID:    userID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	pgRef := &model.AchievementRef{StudentID: userID}
	if err := s.achRepo.Create(ctx, mongoData, pgRef); err != nil {
		return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()})
	}
	return c.Status(201).JSON(helper.SuccessResponse{Success: true, Message: "Achievement Drafted", Data: pgRef})
}

// [FR-005] Update Achievement (Draft Only)
func (s *achService) UpdateAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	var req map[string]interface{} // Flexible body
	c.BodyParser(&req)

	ctx := c.Context()
	data, ref, err := s.achRepo.FindByID(ctx, id)
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "Not Found"}) }

	if ref.StudentID != userID || ref.Status != "draft" {
		return c.Status(403).JSON(helper.ErrorResponse{Success: false, Message: "Cannot edit this achievement"})
	}

	// Update Mongo Data
	if err := s.achRepo.UpdateData(ctx, data.ID.Hex(), req); err != nil {
		return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Updated"})
}

// Delete Draft
func (s *achService) DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	ctx := c.Context()
	_, ref, err := s.achRepo.FindByID(ctx, id)
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "Not Found"}) }

	// Admin bisa hapus apapun, User cuma bisa hapus draft miliknya
	if role != "Admin" {
		if ref.StudentID != userID { return c.Status(403).JSON(helper.ErrorResponse{Success: false, Message: "Forbidden"}) }
		if ref.Status != "draft" { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Only draft can be deleted"}) }
	}

	if err := s.achRepo.DeleteDraft(ctx, id); err != nil {
		return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Deleted"})
}

// Get History
func (s *achService) GetMyHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	data, err := s.achRepo.FindHistoryByUserID(c.Context(), userID)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: data})
}

// Admin Finalize (ACC)
func (s *achService) AdminFinalize(c *fiber.Ctx) error {
	id := c.Params("id")
	adminID := c.Locals("user_id").(string)
	err := s.achRepo.UpdateStatus(c.Context(), id, "finalized", "Approved by Admin", adminID)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: "Failed"}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Finalized by Admin"})
}

// ... (Submit, Verify, Reject, Stats, GetAll Sama seperti sebelumnya, hanya copy paste) ...
func (s *achService) SubmitForVerification(c *fiber.Ctx) error {
	id := c.Params("id")
	err := s.achRepo.UpdateStatus(c.Context(), id, "submitted", "", "")
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Submitted"})
}

func (s *achService) VerifyAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	dosenID := c.Locals("user_id").(string)
	err := s.achRepo.UpdateStatus(c.Context(), id, "verified", "", dosenID)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Verified"})
}

func (s *achService) RejectAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	dosenID := c.Locals("user_id").(string)
	var req model.VerifyRequest
	c.BodyParser(&req)
	err := s.achRepo.UpdateStatus(c.Context(), id, "rejected", req.RejectionNote, dosenID)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Rejected"})
}

func (s *achService) GetAllAchievements(c *fiber.Ctx) error {
	data, err := s.achRepo.FindAll(c.Context())
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: data})
}

func (s *achService) GetStatistics(c *fiber.Ctx) error {
	stats, err := s.achRepo.CountByStatus(c.Context())
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: stats})
}