package repository

import (
	"mmm_server/pkg/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	GetAllUsers() ([]model.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Collection) *Repository {
	return &Repository{
		User: NewUserMongoDb(db),
	}
}
