package model

import "time"

// Auth DTO
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// User Management DTO
type CreateUserRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FullName       string `json:"full_name"`
	RoleName       string `json:"role_name"`
	IdentityNumber string `json:"identity_number"` // NIM/NIP
	Department     string `json:"department"`
	AdvisorID      string `json:"advisor_id"`
}

type UpdateUserRequest struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Password       string `json:"password,omitempty"` // Optional
	IsActive       *bool  `json:"is_active"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Password string `json:"password,omitempty"`
}

type StudentResponse struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	NIM            string `json:"nim"`
	ProgramStudy   string `json:"program_study"`
	AdvisorName    string `json:"advisor_name"` // Nama Dosen Wali
}

// [BARU] Response Khusus Data Dosen Lengkap
type LecturerResponse struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	NIP        string `json:"nip"`
	Department string `json:"department"`
}

// Achievement DTO
type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
}

type VerifyRequest struct {
	Status        string `json:"status"`
	RejectionNote string `json:"rejection_note,omitempty"`
}

// Response Wrapper
type AchievementHistoryResponse struct {
	RefID     string                 `json:"id"`
	Title     string                 `json:"title"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	Details   map[string]interface{} `json:"details,omitempty"`
}