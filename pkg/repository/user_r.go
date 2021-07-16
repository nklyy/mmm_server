package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mmm_server/pkg/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMongoDb struct {
	db *mongo.Database
}

func NewUserMongoDb(db *mongo.Database) *UserMongoDb {
	return &UserMongoDb{db: db}
}

func (ur *UserMongoDb) GetUserMusicDB(guestID string) ([]model.GeneralMusicStruct, error) {
	var music []model.GeneralMusicStruct

	err := ur.db.Collection("user").FindOne(context.TODO(), bson.M{"guest_id": guestID}).Decode(&music)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func (ur *UserMongoDb) CreateGuestUserDB(guestID string) (bool, error) {
	var user model.User
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.GuestId = guestID
	user.Music = []model.GeneralMusicStruct{}

	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(60),
	}

	_, err := ur.db.Collection("user").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return false, err
	}

	_, err = ur.db.Collection("user").InsertOne(context.TODO(), user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ur *UserMongoDb) UpdateGuestUserDB(guestID string, uMusic []model.GeneralMusicStruct) (bool, error) {
	_, err := ur.db.Collection("user").UpdateOne(context.TODO(), bson.M{"guest_id": guestID}, bson.D{{"$set", bson.D{{"music", uMusic}}}})
	if err != nil {
		return false, err
	}

	return true, nil
}
