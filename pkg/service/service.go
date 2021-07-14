package service

import (
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type User interface {
	GetAllUsersDB() ([]model.User, error)
}

type Deezer interface {
	GetDeezerAccessToken(code string) string
	GetDeezerUserMusic(token string) []DeezerTrack
	CheckDeezerAccessToken(token string) bool
}

type Spotify interface {
	GetSpotifyAccessToken(code string) string
	CheckSpotifyAccessToken(token string) bool
	GetSpotifyUserMusic(token string) []SpotifyTrack
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
