package repository

import (
	"context"
	"database/sql"
	"tugas_uas/app/model"
)

type LecturerRepository interface {
	FindAll(ctx context.Context) ([]model.LecturerResponse, error)
	FindByID(ctx context.Context, id string) (*model.LecturerResponse, error)
	GetLecturerIDByUserID(ctx context.Context, userID string) (string, error)
}

type lecturerRepo struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepo{db: db}
}

// Ambil semua dosen (untuk dropdown Admin)
func (r *lecturerRepo) FindAll(ctx context.Context) ([]model.LecturerResponse, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.email, l.lecturer_id, l.department
		FROM lecturers l
		JOIN users u ON l.user_id = u.id
		ORDER BY u.full_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var lecturers []model.LecturerResponse
	for rows.Next() {
		var l model.LecturerResponse
		rows.Scan(&l.UserID, &l.Username, &l.FullName, &l.Email, &l.NIP, &l.Department)
		lecturers = append(lecturers, l)
	}
	return lecturers, nil
}

// Cari Detail Dosen by ID
func (r *lecturerRepo) FindByID(ctx context.Context, id string) (*model.LecturerResponse, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.email, l.lecturer_id, l.department
		FROM lecturers l
		JOIN users u ON l.user_id = u.id
		WHERE l.id = $1
	`
	var l model.LecturerResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(&l.UserID, &l.Username, &l.FullName, &l.Email, &l.NIP, &l.Department)
	return &l, err
}

// Helper: Ubah UserID (Login) jadi LecturerID (Tabel)
func (r *lecturerRepo) GetLecturerIDByUserID(ctx context.Context, userID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM lecturers WHERE user_id = $1", userID).Scan(&id)
	return id, err
}