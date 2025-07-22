package config

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func Connection() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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
    MongoDatabase = client.Database("sateayammadura")
}