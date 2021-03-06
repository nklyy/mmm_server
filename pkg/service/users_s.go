package service

import (
	"mmm_server/config"
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
)

type UserService struct {
	repo repository.User
	cfg  *config.Configurations
}

func NewUserService(repo repository.User, cfg *config.Configurations) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (us *UserService) GetUser(guestID string) (model.User, error) {
	user, err := us.repo.GetGuestUserDB(guestID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *UserService) CreateGuestUser(guestID string, findAccessT string) error {
	err := us.repo.CreateGuestUserDB(guestID, findAccessT)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdateGuestUser(guestID string, user model.User) error {
	err := us.repo.UpdateGuestUserDB(guestID, user)
	if err != nil {
		return err
	}

	return nil
}
