package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeneralMusicStruct struct {
	ID       string
	SongName string
	//ArtistName string
	AlbumName string
}

type User struct {
	ID        primitive.ObjectID   `bson:"_id"`
	Music     []GeneralMusicStruct `bson:"music"`
	GuestId   string               `bson:"guest_id"`
	CreatedAt time.Time            `bson:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at"`
}
