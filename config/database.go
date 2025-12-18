package config

import (
	"database/sql"
	"tugas_uas/database"

	"go.mongodb.org/mongo-driver/mongo"
)

// Container untuk menyimpan kedua koneksi
type DatabaseContainer struct {
	Postgres *sql.DB
	Mongo    *mongo.Database
}

func InitDB() *DatabaseContainer {
	return &DatabaseContainer{
		Postgres: database.ConnectPostgres(),
		Mongo:    database.ConnectMongo(),
	}
}