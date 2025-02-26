package repositories

import (
	"errors"

	"github.com/domnikl/music-box-game/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) CreateUser(user *models.User) error {
	err := u.db.Create(user)
	if err != nil {
		return err.Error
	}

	return nil
}

func (u *UserRepository) UpdateUser(user *models.User) error {
	err := u.db.Save(user)
	if err != nil {
		return err.Error
	}

	return nil
}

func (u *UserRepository) FindUserByAPIToken(apiToken string) (*models.User, error) {
	var user models.User
	err := u.db.Where("api_token = ?", apiToken).First(&user)
	if err != nil && err.Error != nil {
		return nil, err.Error
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
