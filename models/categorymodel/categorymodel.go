package categorymodel

import (
    "context"
    "sateayammadura/config"
    "sateayammadura/entities"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
    

)

func GetAllCategories() ([]entities.Category, error) {
    var categories []entities.Category
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    collection := config.MongoDatabase.Collection("categories")
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var category entities.Category
        if err := cursor.Decode(&category); err != nil {
            continue
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func Create(category entities.Category) bool {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.MongoDatabase.Collection("categories")
	_, err := collection.InsertOne(ctx, category)
	return err == nil
}

func Detail(id primitive.ObjectID) entities.Category {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var category entities.Category
    err := config.MongoDatabase.
        Collection("categories").
        FindOne(ctx, bson.M{"_id": id}).
        Decode(&category)

    if err != nil {
        return entities.Category{}
    }

    return category
}

func Update(id primitive.ObjectID, category entities.Category) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    category.UpdatedAt = time.Now()
    _, err := config.MongoDatabase.
        Collection("categories").
        UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
            "name":       category.Name,
            "updated_at": category.UpdatedAt,
        }})

    return err == nil
}

func Delete(id primitive.ObjectID) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := config.MongoDatabase.
        Collection("categories").
        DeleteOne(ctx, bson.M{"_id": id})

    return err == nil
}

func GetCategoryByID(id primitive.ObjectID) (entities.Category, error) {
	var category entities.Category
	collection := config.MongoDatabase.Collection("categories")
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&category)
	return category, err
}


