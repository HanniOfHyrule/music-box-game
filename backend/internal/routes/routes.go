package routes

import (
	"github.com/domnikl/music-box-game/backend/internal/spotify"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Setup(e *echo.Echo, db *gorm.DB, spotify *spotify.Spotify) {
	setupAuth(e, db)
	setupSpotify(e, db, spotify)
}
