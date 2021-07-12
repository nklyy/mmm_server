package service

import (
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type User interface {
	GetAllUsers() ([]model.User, error)
}

type Deezer interface {
	GetDeezerAccessToken(code string) string
	GetDeezerUserMusic(token string) []model.Track
	CheckAccessToken(token string) bool
}

type Service struct {
	User
	Deezer
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:   NewUserService(repos.User),
		Deezer: NewDeezerService(repos.User),
	}
}
