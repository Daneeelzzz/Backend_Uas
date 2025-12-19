package service

import (
	"context"
	"math" // Untuk hitung total page
	"strconv" // Untuk convert string ke int
	"strings" // Untuk upper case order
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

// [FR-009] Admin Get All Users with Pagination
// @Summary Lihat semua user (Paginated)
// @Description Mengambil daftar user dengan opsi page, limit, sort_by, dan order
// @Tags Admin User Management
// @Param page query int false "Nomor Halaman (Default 1)"
// @Param limit query int false "Jumlah data per halaman (Default 10)"
// @Param sort_by query string false "Kolom sort (username, full_name, email)"
// @Param order query string false "Urutan (asc, desc)"
// @Success 200 {object} helper.SuccessResponse{data=[]model.User}
// @Router /admin/users [get]
func (s *userService) GetAllUsers(c *fiber.Ctx) error {
	// 1. Ambil Parameter (Query Params)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "created_at")
	order := strings.ToUpper(c.Query("order", "DESC"))

	// Validasi dasar
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }

	// Hitung Offset
	offset := (page - 1) * limit

	ctx := c.Context()
	
	// 2. Panggil Repo
	users, totalData, err := s.repo.FindAll(ctx, limit, offset, sortBy, order)
	if err != nil {
		return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()})
	}

	// 3. Hitung Metadata
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	meta := &helper.Meta{
		CurrentPage: page,
		TotalPage:   totalPage,
		TotalData:   totalData,
		Limit:       limit,
	}

	// 4. Return Response dengan Meta
	return c.JSON(helper.SuccessResponse{
		Success: true, 
		Message: "Success Fetch Users", 
		Data: users,
		Meta: meta, // Sertakan Meta info
	})
}

// ... (Method CRUD lainnya CreateUser, UpdateUser, dll SAMA SEPERTI SEBELUMNYA) ...
// Saya skip tulis ulang untuk menghemat space, tapi pastikan kamu tidak menghapus method lain di file ini.

func (s *userService) CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Invalid Request"}) }
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second); defer cancel()
	hashed, _ := utils.HashPassword(req.Password)
	roleID, err := s.repo.GetRoleIDByName(ctx, req.RoleName)
	if err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Role invalid"}) }
	user := &model.User{Username: req.Username, Email: req.Email, PasswordHash: hashed, FullName: req.FullName, RoleID: roleID}
	extra := map[string]string{"identity_number": req.IdentityNumber, "department": req.Department, "advisor_id": req.AdvisorID}
	if err := s.repo.CreateUser(ctx, user, req.RoleName, extra); err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false, Message: err.Error()}) }
	return c.Status(201).JSON(helper.SuccessResponse{Success: true, Message: "User Created", Data: user})
}

func (s *userService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id"); var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false}) }
	ctx := c.Context()
	user, err := s.repo.FindByID(ctx, id)
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "User Not Found"}) }
	if req.FullName != "" { user.FullName = req.FullName }
	if req.Email != "" { user.Email = req.Email }
	if req.Password != "" { hash, _ := utils.HashPassword(req.Password); user.PasswordHash = hash }
	if req.IsActive != nil { user.IsActive = *req.IsActive }
	if err := s.repo.UpdateUser(ctx, id, user); err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "User Updated"})
}

func (s *userService) DeleteUser(c *fiber.Ctx) error {
	if err := s.repo.DeleteUser(c.Context(), c.Params("id")); err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "User Deleted"})
}

func (s *userService) GetMyProfile(c *fiber.Ctx) error {
	user, err := s.repo.FindByID(c.Context(), c.Locals("user_id").(string))
	if err != nil { return c.Status(404).JSON(helper.ErrorResponse{Success: false, Message: "Profile not found"}) }
	return c.JSON(helper.SuccessResponse{Success: true, Data: user})
}

func (s *userService) UpdateMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string); var req model.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(helper.ErrorResponse{Success: false}) }
	ctx := c.Context()
	user, _ := s.repo.FindByID(ctx, userID)
	if req.FullName != "" { user.FullName = req.FullName }
	if req.Password != "" { hash, _ := utils.HashPassword(req.Password); user.PasswordHash = hash }
	if err := s.repo.UpdateUser(ctx, userID, user); err != nil { return c.Status(500).JSON(helper.ErrorResponse{Success: false}) }
	return c.JSON(helper.SuccessResponse{Success: true, Message: "Profile Updated"})
}