package repository

import (
	"context"
	"database/sql"
	"tugas_uas/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(ctx context.Context, data *model.AchievementMongo, ref *model.AchievementRef) error
	UpdateStatus(ctx context.Context, refID string, status string, notes string, verifierID string) error
	CountByStatus(ctx context.Context) (map[string]int, error)
	FindHistoryByUserID(ctx context.Context, userID string) ([]model.AchievementRef, error)
	FindAll(ctx context.Context) ([]model.AchievementRef, error)
	FindByID(ctx context.Context, refID string) (*model.AchievementMongo, *model.AchievementRef, error)
	DeleteDraft(ctx context.Context, refID string) error
	
	// NEW: Update Data Mongo (Edit Draft)
	UpdateData(ctx context.Context, mongoID string, data map[string]interface{}) error
}

type achievementRepo struct {
	pg    *sql.DB
	mongo *mongo.Collection
}

func NewAchievementRepository(pg *sql.DB, mongoDB *mongo.Database) AchievementRepository {
	return &achievementRepo{pg: pg, mongo: mongoDB.Collection("achievements")}
}

// ... (Method Create, UpdateStatus, CountByStatus, FindHistoryByUserID, FindAll, FindByID SAMA SEPERTI SEBELUMNYA) ...
// Copy dari jawaban sebelumnya untuk method yang tidak berubah agar hemat space.
// Saya tulis ulang method Create dan FindByID agar konteksnya jelas.

func (r *achievementRepo) Create(ctx context.Context, data *model.AchievementMongo, ref *model.AchievementRef) error {
	res, err := r.mongo.InsertOne(ctx, data)
	if err != nil { return err }
	mongoID := res.InsertedID.(primitive.ObjectID).Hex()
	ref.MongoAchievementID = mongoID
	query := `INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, created_at) VALUES (gen_random_uuid(), $1, $2, 'draft', NOW()) RETURNING id`
	err = r.pg.QueryRowContext(ctx, query, ref.StudentID, ref.MongoAchievementID).Scan(&ref.ID)
	if err != nil { r.mongo.DeleteOne(ctx, bson.M{"_id": res.InsertedID}); return err }
	return nil
}

func (r *achievementRepo) UpdateStatus(ctx context.Context, refID string, status string, notes string, verifierID string) error {
	query := `UPDATE achievement_references SET status = $1, rejection_note = $2, verified_by = $3, updated_at = NOW() WHERE id = $4`
	var vID interface{} = nil
	if verifierID != "" { vID = verifierID }
	_, err := r.pg.ExecContext(ctx, query, status, notes, vID, refID)
	return err
}

func (r *achievementRepo) CountByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := r.pg.QueryContext(ctx, `SELECT status, COUNT(*) FROM achievement_references GROUP BY status`)
	if err != nil { return nil, err }
	defer rows.Close()
	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		stats[status] = count
	}
	return stats, nil
}

func (r *achievementRepo) FindHistoryByUserID(ctx context.Context, userID string) ([]model.AchievementRef, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, created_at FROM achievement_references WHERE student_id = $1 ORDER BY created_at DESC`
	rows, err := r.pg.QueryContext(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var results []model.AchievementRef
	for rows.Next() {
		var a model.AchievementRef
		rows.Scan(&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.CreatedAt)
		results = append(results, a)
	}
	return results, nil
}

func (r *achievementRepo) FindAll(ctx context.Context) ([]model.AchievementRef, error) {
	query := `SELECT id, student_id, mongo_achievement_id, status, created_at FROM achievement_references ORDER BY created_at DESC`
	rows, err := r.pg.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()
	var results []model.AchievementRef
	for rows.Next() {
		var a model.AchievementRef
		rows.Scan(&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.CreatedAt)
		results = append(results, a)
	}
	return results, nil
}

func (r *achievementRepo) FindByID(ctx context.Context, refID string) (*model.AchievementMongo, *model.AchievementRef, error) {
	var ref model.AchievementRef
	query := `SELECT id, student_id, mongo_achievement_id, status FROM achievement_references WHERE id = $1`
	err := r.pg.QueryRowContext(ctx, query, refID).Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status)
	if err != nil { return nil, nil, err }
	objID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	var data model.AchievementMongo
	err = r.mongo.FindOne(ctx, bson.M{"_id": objID}).Decode(&data)
	return &data, &ref, err
}

func (r *achievementRepo) DeleteDraft(ctx context.Context, refID string) error {
	var mongoID string
	err := r.pg.QueryRowContext(ctx, "SELECT mongo_achievement_id FROM achievement_references WHERE id = $1 AND status = 'draft'", refID).Scan(&mongoID)
	if err != nil { return err }
	_, err = r.pg.ExecContext(ctx, "DELETE FROM achievement_references WHERE id = $1", refID)
	if err != nil { return err }
	objID, _ := primitive.ObjectIDFromHex(mongoID)
	_, err = r.mongo.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// [NEW] Implementasi Edit Data Mongo
func (r *achievementRepo) UpdateData(ctx context.Context, mongoID string, data map[string]interface{}) error {
	objID, _ := primitive.ObjectIDFromHex(mongoID)
	update := bson.M{"$set": data}
	_, err := r.mongo.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}