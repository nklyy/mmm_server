package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"mmm_server/pkg/model"
)

type User interface {
	GetUserDB(guestID string) (model.User, error)
	CreateGuestUserDB(guestID string, findAccessToken string) (bool, error)
	UpdateGuestUserDB(guestID string, user model.User) (bool, error)
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User: NewUserMongoDb(db),
	}
}
