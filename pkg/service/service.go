package service

import (
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type User interface {
	GetUserMusic(guestID string) ([]model.GeneralMusicStruct, error)
	CreateGuestUser(guestID string, accessT string)
	UpdateGuestUser(guestID string, uMusic []model.GeneralMusicStruct)
}

type Deezer interface {
	GetDeezerAccessToken(code string) string
	GetDeezerUserMusic(guestID string) []model.GeneralMusicStruct
	CheckDeezerAccessToken(guestID string) bool
	MoveToDeezer(tracks []model.GeneralMusicStruct, guestID string)
}

type Spotify interface {
	GetSpotifyAccessToken(code string) string
	GetSpotifyUserMusic(guestID string) []model.GeneralMusicStruct
	CheckSpotifyAccessToken(guestID string) bool
	MoveToSpotify(tracks []model.GeneralMusicStruct, guestID string)
}

type Service struct {
	User
	Deezer
	Spotify
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:    NewUserService(repos.User),
		Deezer:  NewDeezerService(repos.User),
		Spotify: NewSpotifyService(repos.User),
	}
}
