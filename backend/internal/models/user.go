package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint           `json:"id" gorm:"primarykey"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
	APIToken            string         `json:"apiToken"`
	SpotifyToken        string         `json:"-"`
	SpotifyRefreshToken string         `json:"-"`
}
