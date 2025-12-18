package service

import (
	"context"
	"time"
	"tugas_uas/app/model"
	"tugas_uas/app/repository"
	"tugas_uas/helper"
	"tugas_uas/utils"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	CreateUser(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetMyProfile(c *fiber.Ctx) error
	UpdateMyProfile(c *fiber.Ctx) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser godoc
// @Summary Admin membuat user baru
// @Description Membuat user baru dengan role tertentu (Admin/Mahasiswa/Dosen)
// @Tags Admin User Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body model.CreateUserRequest true "Data User Baru"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Failure 500 {object} helper.ErrorResponse
// @Security BearerAuth
// @Router /admin/users [post]
func (s *userService) CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Invalid Request"}) }
	
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	hashed, err := utils.HashPassword(req.Password)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: "Hashing Failed"}) }

	roleID, err := s.repo.GetRoleIDByName(ctx, req.RoleName)
	if err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Role invalid"}) }

	user := &model.User{Username: req.Username, Email: req.Email, PasswordHash: hashed, FullName: req.FullName, RoleID: roleID}
	extra := map[string]string{"identity_number": req.IdentityNumber, "department": req.Department, "advisor_id": req.AdvisorID}

	err = s.repo.CreateUser(ctx, user, req.RoleName, extra)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }

	return c.Status(201).JSON(helper.SuccessResponse{Success: true, Message: "User Created", Data: user})
}

// GetAllUsers godoc
// @Summary Lihat semua user
// @Description Mengambil daftar semua user di sistem
// @Tags Admin User Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} helper.SuccessResponse{data=[]model.User}
// @Security BearerAuth
// @Router /admin/users [get]
func (s *userService) GetAllUsers(c *fiber.Ctx) error {
	users, err := s.repo.FindAll(c.Context())
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: users})
}

// ... (Method Update, Delete, dll biarkan seperti file sebelumnya, tambahkan comment jika mau) ...

func (s *userService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Bad Input"}) }

	ctx := c.Context()
	user, err := s.repo.FindByID(ctx, id)
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "User Not Found"}) }

	if req.FullName != "" { user.FullName = req.FullName }
	if req.Email != "" { user.Email = req.Email }
	if req.Password != "" { 
		hash, _ := utils.HashPassword(req.Password)
		user.PasswordHash = hash 
	}
	if req.IsActive != nil { user.IsActive = *req.IsActive }

	err = s.repo.UpdateUser(ctx, id, user)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "User Updated"})
}

func (s *userService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := s.repo.DeleteUser(c.Context(), id)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "User Deleted"})
}

func (s *userService) GetMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := s.repo.FindByID(c.Context(), userID)
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "Profile not found"}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: user})
}

func (s *userService) UpdateMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req model.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Bad Input"}) }
	
	ctx := c.Context()
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false}) }

	if req.FullName != "" { user.FullName = req.FullName }
	if req.Password != "" { 
		hash, _ := utils.HashPassword(req.Password)
		user.PasswordHash = hash 
	}

	err = s.repo.UpdateUser(ctx, userID, user)
	if err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Profile Updated"})
}