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

type AuthService interface {
	Login(c *fiber.Ctx) error
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

// LoginHandler godoc
// @Summary Masuk ke dalam sistem
// @Description Authenticate user dengan username/email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Kredensial Login"
// @Success 200 {object} helper.SuccessResponse{data=model.LoginResponse}
// @Failure 400 {object} helper.ErrorResponse
// @Failure 401 {object} helper.ErrorResponse
// @Router /auth/login [post]
func (s *authService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(helper.ErrorResponse{Success: false, Message: "Invalid Request"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	user, err := s.repo.FindByEmailOrUsername(ctx, req.Identifier)
	if err != nil {
		return c.Status(401).JSON(helper.ErrorResponse{Success: false, Message: "User not found"})
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(helper.ErrorResponse{Success: false, Message: "Wrong password"})
	}

	if !user.IsActive {
		return c.Status(403).JSON(helper.ErrorResponse{Success: false, Message: "Account is inactive"})
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email, user.RoleName)

	return c.JSON(helper.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: model.LoginResponse{
			Token: token,
			User:  *user,
		},
	})
}