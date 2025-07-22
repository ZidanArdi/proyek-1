// entities/category.go
package entities

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Category struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`  // ⬅️ ganti dari interface{} jadi primitive.ObjectID
    Name      string             `bson:"name"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}
