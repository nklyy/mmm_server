package service

import (
	"github.com/gofiber/websocket/v2"
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type User interface {
	GetUser(guestID string) (model.User, error)
	CreateGuestUser(guestID string, findAccessToken string)
	UpdateGuestUser(guestID string, user model.User)
}

type Deezer interface {
	GetDeezerAccessToken(code string) string
	GetDeezerUserMusic(guestID string) []model.GeneralMusicStruct
	CheckDeezerAccessToken(guestID string) bool
	MoveToDeezer(accessToken string, tracks []model.GeneralMusicStruct, con *websocket.Conn, mt int)
}

type Spotify interface {
	GetSpotifyAccessToken(code string) string
	GetSpotifyUserMusic(guestID string) []model.GeneralMusicStruct
	CheckSpotifyAccessToken(guestID string) bool
	MoveToSpotify(accessToken string, tracks []model.GeneralMusicStruct, con *websocket.Conn, mt int)
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
