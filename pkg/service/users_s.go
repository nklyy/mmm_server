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

func (us *UserService) GetUserMusic(guestID string) ([]model.GeneralMusicStruct, error) {
	music, err := us.repo.GetUserMusicDB(guestID)

	if err != nil {
		return nil, err
	}

	return music, nil
}

func (us *UserService) CreateGuestUser(guestID string) {
	us.repo.CreateGuestUserDB(guestID)
}

func (us *UserService) UpdateGuestUser(guestID string, uMusic []model.GeneralMusicStruct) {
	us.repo.UpdateGuestUserDB(guestID, uMusic)
}
