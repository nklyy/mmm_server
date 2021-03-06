package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeneralMusicStruct struct {
	ID         string
	ArtistName string
	SongName   string
	AlbumName  string
}

type User struct {
	ID              primitive.ObjectID   `bson:"_id"`
	Music           []GeneralMusicStruct `bson:"music"`
	GuestId         string               `bson:"guest_id"`
	AccessTokenFind string               `bson:"access_token_find"`
	AccessTokenMove string               `bson:"access_token_move"`
	CreatedAt       time.Time            `bson:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at"`
}
