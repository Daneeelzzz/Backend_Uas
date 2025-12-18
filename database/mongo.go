package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ Gagal inisialisasi MongoDB:", err)
	}

	// Ping untuk memastikan koneksi hidup
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Gagal koneksi MongoDB:", err)
	}

	log.Println("✅ MongoDB Connected")
	return client.Database(os.Getenv("MONGO_DB_NAME"))
}