package middlewares

import (
	"net/http"

	"github.com/domnikl/music-box-game/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type SpotifyMiddleware struct {
}

func NewSpotifyMiddleware() SpotifyMiddleware {
	return SpotifyMiddleware{}
}

func (m SpotifyMiddleware) HasToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Unauthorized",
			})
		}

		if user.SpotifyRefreshToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing Spotify token",
			})
		}

		return next(c)
	}
}
