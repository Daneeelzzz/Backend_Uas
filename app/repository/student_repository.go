package repository

import (
	"context"
	"database/sql"
	"tugas_uas/app/model"
)

type StudentRepository interface {
	FindAll(ctx context.Context) ([]model.StudentResponse, error)
	
	// [RELASI ONE-TO-MANY] Cari mahasiswa berdasarkan ID Dosen Walinya
	FindByAdvisorID(ctx context.Context, advisorID string) ([]model.StudentResponse, error)
	
	UpdateAdvisor(ctx context.Context, studentUserID string, advisorID string) error
}

type studentRepo struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepo{db: db}
}

// Ambil semua mahasiswa (Admin View)
func (r *studentRepo) FindAll(ctx context.Context) ([]model.StudentResponse, error) {
	query := `
		SELECT 
			u.id, u.username, u.full_name, u.email, 
			s.student_id, s.program_study,
			COALESCE(adv_u.full_name, '-') as advisor_name
		FROM students s
		JOIN users u ON s.user_id = u.id
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
		rows.Scan(&s.UserID, &s.Username, &s.FullName, &s.Email, &s.NIM, &s.ProgramStudy, &s.AdvisorName)
		students = append(students, s)
	}
	return students, nil
}

// [CORE ONE-TO-MANY] Ambil semua anak wali dari satu dosen
func (r *studentRepo) FindByAdvisorID(ctx context.Context, advisorID string) ([]model.StudentResponse, error) {
	query := `
		SELECT 
			u.id, u.username, u.full_name, u.email, 
			s.student_id, s.program_study,
			'ME' as advisor_name
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.advisor_id = $1 -- INI KUNCINYA
		ORDER BY u.full_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, advisorID)
	if err != nil { return nil, err }
	defer rows.Close()

	var students []model.StudentResponse
	for rows.Next() {
		var s model.StudentResponse
		rows.Scan(&s.UserID, &s.Username, &s.FullName, &s.Email, &s.NIM, &s.ProgramStudy, &s.AdvisorName)
		students = append(students, s)
	}
	return students, nil
}

// Admin Assign Dosen Wali ke Mahasiswa
func (r *studentRepo) UpdateAdvisor(ctx context.Context, studentUserID string, advisorID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE students SET advisor_id = $1 WHERE user_id = $2`, advisorID, studentUserID)
	return err
}