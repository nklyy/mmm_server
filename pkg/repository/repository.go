package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"mmm_server/pkg/model"
)

type User interface {
	GetGuestUserDB(guestID string) (model.User, error)
	CreateGuestUserDB(guestID string, findAccessToken string) error
	UpdateGuestUserDB(guestID string, user model.User) error
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User: NewUserMongoDb(db),
	}
}
