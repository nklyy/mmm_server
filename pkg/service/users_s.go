package service

import (
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetAllUsersDB() ([]model.User, error) {
	return us.repo.GetAllUsers("")
}
