package repositories

import (
	"mmm_server/repositories/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	GetAllUsers() []models.User
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Collection) *Repository {
	return &Repository{
		User: NewUserMongoDb(db),
	}
}
