package repository

import (
	"context"
	"log"
	"mmm_server/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoDb struct {
	db *mongo.Database
}

func NewUserMongoDb(db *mongo.Database) *UserMongoDb {
	return &UserMongoDb{db: db}
}

func (ur *UserMongoDb) GetAllUsers(sort string) ([]model.User, error) {
	var users []model.User

	cursor, err := ur.db.Collection("user").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	err = cursor.All(context.TODO(), &users)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(users)

	return users, err
}
