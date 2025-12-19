package repository

import (
	"context"
	"database/sql"
	"fmt"
	"errors"
	"tugas_uas/app/model"
)

type UserRepository interface {
	FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	GetRoleIDByName(ctx context.Context, name string) (string, error)
	
	CreateUser(ctx context.Context, user *model.User, roleName string, extraData map[string]string) error
	FindAll(ctx context.Context, limit, offset int, sortBy, order string) ([]model.User, int, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	
	// Method ini masih disimpan di UserRepo karena sering dipakai umum,
	// tapi implementasi logic utamanya ada di StudentRepo/LecturerRepo untuk query spesifik.
	// Jika method GetStudentsByAdvisor sudah tidak dipakai di user_service, boleh dihapus dari interface ini.
	// Namun untuk amannya kita biarkan dulu jika ada logic lain yang pakai.
	GetStudentsByAdvisor(ctx context.Context, lecturerUserID string) ([]model.User, error)
	
	FindAllStudents(ctx context.Context) ([]model.StudentResponse, error)
	FindAllLecturers(ctx context.Context) ([]model.LecturerResponse, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	query := `SELECT u.id, u.username, u.email, u.password_hash, u.full_name, r.name, u.is_active FROM users u JOIN roles r ON u.role_id = r.id WHERE u.email = $1 OR u.username = $1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.RoleName, &user.IsActive)
	return &user, err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT u.id, u.username, u.email, u.password_hash, u.full_name, r.name, u.is_active FROM users u JOIN roles r ON u.role_id = r.id WHERE u.id = $1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.RoleName, &user.IsActive)
	return &user, err
}

func (r *userRepository) GetRoleIDByName(ctx context.Context, name string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM roles WHERE name = $1", name).Scan(&id)
	if err == sql.ErrNoRows {
		return "", errors.New("role not found")
	}
	return id, err
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User, roleName string, extraData map[string]string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil { return err }

	queryUser := `INSERT INTO users (username, email, password_hash, full_name, role_id, is_active) VALUES ($1, $2, $3, $4, $5, true) RETURNING id`
	err = tx.QueryRowContext(ctx, queryUser, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID).Scan(&user.ID)
	if err != nil { tx.Rollback(); return err }

	// [FIX] Menggunakan switch case agar lebih clean dan lolos linter
	switch roleName {
	case "Mahasiswa":
		var advisorID interface{}
		if val, ok := extraData["advisor_id"]; ok && val != "" { advisorID = val } else { advisorID = nil }
		_, err = tx.ExecContext(ctx, `INSERT INTO students (user_id, student_id, program_study, advisor_id) VALUES ($1, $2, $3, $4)`, user.ID, extraData["identity_number"], extraData["department"], advisorID)
	case "Dosen Wali":
		_, err = tx.ExecContext(ctx, `INSERT INTO lecturers (user_id, lecturer_id, department) VALUES ($1, $2, $3)`, user.ID, extraData["identity_number"], extraData["department"])
	}

	if err != nil { tx.Rollback(); return err }
	return tx.Commit()
}

// [UPDATE] FindAll dengan Pagination & Sorting
func (r *userRepository) FindAll(ctx context.Context, limit, offset int, sortBy, order string) ([]model.User, int, error) {
	// 1. Whitelist sorting columns (Security Best Practice: SQL Injection Prevention)
	validColumns := map[string]bool{"username": true, "email": true, "full_name": true, "created_at": true}
	if !validColumns[sortBy] {
		sortBy = "created_at" // Default
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC" // Default
	}

	// 2. Query Utama dengan Pagination
	query := fmt.Sprintf(`
		SELECT u.id, u.username, u.email, u.full_name, r.name, u.is_active 
		FROM users u 
		JOIN roles r ON u.role_id = r.id 
		ORDER BY u.%s %s 
		LIMIT $1 OFFSET $2`, sortBy, order)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleName, &u.IsActive); err == nil {
			users = append(users, u)
		}
	}

	// 3. Hitung Total Data (Untuk Metadata Pagination)
	var total int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id string, user *model.User) error {
	query := `UPDATE users SET full_name = $1, email = $2, password_hash = $3, is_active = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.FullName, user.Email, user.PasswordHash, user.IsActive, id)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepository) GetStudentsByAdvisor(ctx context.Context, lecturerUserID string) ([]model.User, error) {
	query := `SELECT u.id, u.username, u.full_name, u.email FROM users u JOIN students s ON u.id = s.user_id JOIN lecturers l ON s.advisor_id = l.id WHERE l.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, lecturerUserID)
	if err != nil { return nil, err }
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Username, &u.FullName, &u.Email)
		u.RoleName = "Mahasiswa"
		users = append(users, u)
	}
	return users, nil
}

// Method untuk support StudentRepo & LecturerRepo yang mungkin memanggil fungsi ini via embedding atau logic terpisah
func (r *userRepository) FindAllStudents(ctx context.Context) ([]model.StudentResponse, error) {
	query := `
		SELECT 
			u.id, u.username, u.full_name, u.email, 
			s.student_id, s.program_study,
			COALESCE(adv_u.full_name, '-') as advisor_name
		FROM users u
		JOIN students s ON u.id = s.user_id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		LEFT JOIN users adv_u ON l.user_id = adv_u.id
		ORDER BY u.full_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var students []model.StudentResponse
	for rows.Next() {
		var s model.StudentResponse
		if err := rows.Scan(&s.UserID, &s.Username, &s.FullName, &s.Email, &s.NIM, &s.ProgramStudy, &s.AdvisorName); err == nil {
			students = append(students, s)
		}
	}
	return students, nil
}

func (r *userRepository) FindAllLecturers(ctx context.Context) ([]model.LecturerResponse, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.email, l.lecturer_id, l.department
		FROM users u
		JOIN lecturers l ON u.id = l.user_id
		ORDER BY u.full_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var lecturers []model.LecturerResponse
	for rows.Next() {
		var l model.LecturerResponse
		if err := rows.Scan(&l.UserID, &l.Username, &l.FullName, &l.Email, &l.NIP, &l.Department); err == nil {
			lecturers = append(lecturers, l)
		}
	}
	return lecturers, nil
}