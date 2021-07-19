package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"mmm_server/pkg/model"
)

type User interface {
	GetUserInfo(guestID string) (model.User, error)
	GetUserMusicDB(guestID string) ([]model.GeneralMusicStruct, error)
	CreateGuestUserDB(guestID string, accessT string) (bool, error)
	UpdateGuestUserDB(guestID string, uMusic []model.GeneralMusicStruct) (bool, error)
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User: NewUserMongoDb(db),
	}
}
