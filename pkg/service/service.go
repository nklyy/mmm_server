package service

import (
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type User interface {
	GetUserMusic(guestID string) ([]model.GeneralMusicStruct, error)
	CreateGuestUser(guestID string)
	UpdateGuestUser(guestID string, uMusic []model.GeneralMusicStruct)
}

type Deezer interface {
	GetDeezerAccessToken(code string) string
	GetDeezerUserMusic(code string) []model.GeneralMusicStruct
	CheckDeezerAccessToken(code string) bool
	MoveToDeezer(tracks []model.GeneralMusicStruct, code string)
}

type Spotify interface {
	GetSpotifyAccessToken(code string) string
	CheckSpotifyAccessToken(code string) bool
	GetSpotifyUserMusic(code string) []model.GeneralMusicStruct
	MoveToSpotify(tracks []model.GeneralMusicStruct, code string)
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
