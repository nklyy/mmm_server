package service

import (
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type User interface {
	GetAllUsers() ([]model.User, error)
}

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}
