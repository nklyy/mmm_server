package repositories

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoDb struct {
	db *mongo.Collection
}

func NewUserMongoDb(db *mongo.Collection) *UserMongoDb {
	return &UserMongoDb{db: db}
}

type U struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Age       int                `bson:"age"`
}

func (ur *UserMongoDb) GetAllUsers() []U {
	cursor, err := ur.db.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var users []U
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(users)

	return users
}
