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

func (us *UserService) GetUser(guestID string) (model.User, error) {
	user, err := us.repo.GetUserDB(guestID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *UserService) CreateGuestUser(guestID string, findAccessT string) {
	_, err := us.repo.CreateGuestUserDB(guestID, findAccessT)
	if err != nil {
		return
	}
}

func (us *UserService) UpdateGuestUser(guestID string, user model.User) {
	_, err := us.repo.UpdateGuestUserDB(guestID, user)
	if err != nil {
		return
	}
}
