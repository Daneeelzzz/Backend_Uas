package model

// [SRS 3.1.1] Tabel users (Postgres)
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	RoleName     string    `json:"role"` // Field hasil JOIN
	IsActive     bool      `json:"is_active"`
}

type Student struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StudentID string `json:"student_id"`
	AdvisorID string `json:"advisor_id"`
}