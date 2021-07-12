package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Music     []string           `bson:"music"`
	DeezerID  string             `bson:"deezer_id"`
	SpotifyID string             `bson:"spotify_id"`
}
