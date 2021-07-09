package repositories

import (
	"context"
	"log"
	"mmm_server/repositories/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoDb struct {
	db *mongo.Collection
}

func NewUserMongoDb(db *mongo.Collection) *UserMongoDb {
	return &UserMongoDb{db: db}
}

func (ur *UserMongoDb) GetAllUsers() []models.User {
	cursor, err := ur.db.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var users []models.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(users)

	return users
}
