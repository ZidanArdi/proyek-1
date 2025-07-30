package config

import (
	"context"
	"log"
	"time"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func Connection() {
	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set in environment")
	}
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		log.Fatal("MONGODB_DATABASE not set in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping untuk memastikan koneksi berhasil
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
	MongoClient = client
	MongoDatabase = client.Database(dbName)
}
