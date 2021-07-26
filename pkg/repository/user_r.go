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

func (ur *UserMongoDb) GetGuestUserDB(guestID string) (model.User, error) {
	var user model.User

	err := ur.db.Collection("user").FindOne(context.TODO(), bson.M{"guest_id": guestID}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserMongoDb) CreateGuestUserDB(guestID string, findAccessT string) error {
	var user model.User
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.GuestId = guestID
	user.Music = []model.GeneralMusicStruct{}
	user.AccessTokenFind = findAccessT

	mod := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetExpireAfterSeconds(60),
	}

	_, err := ur.db.Collection("user").Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return err
	}

	_, err = ur.db.Collection("user").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserMongoDb) UpdateGuestUserDB(guestID string, user model.User) error {
	_, err := ur.db.Collection("user").UpdateOne(context.TODO(), bson.M{"guest_id": guestID}, bson.D{{"$set", user}})
	if err != nil {
		return err
	}

	return nil
}
