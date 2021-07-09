package repositories

import "go.mongodb.org/mongo-driver/mongo"

type User interface {
	GetAllUsers() []U
}

type Repository struct {
	User
}

func NewRepository(db *mongo.Collection) *Repository {
	return &Repository{
		User: NewUserMongoDb(db),
	}
}
