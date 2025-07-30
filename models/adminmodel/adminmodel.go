package adminmodel

import (
	"context"
	"sateayammadura/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func GetAdminByUsername(username string) (*Admin, error) {
	collection := config.MongoDatabase.Collection("admins")
	var admin Admin
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}
