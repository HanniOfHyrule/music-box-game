package services

import (
	"github.com/domnikl/music-box-game/backend/internal/models"
	"github.com/domnikl/music-box-game/backend/internal/repositories"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{repository}
}

func (u *UserService) CreateUser(user *models.User) error {
	return u.repository.CreateUser(user)
}

func (u *UserService) UpdateUser(user *models.User) error {
	return u.repository.UpdateUser(user)
}
