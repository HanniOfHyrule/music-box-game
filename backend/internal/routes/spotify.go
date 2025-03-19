package routes

import (
	"github.com/domnikl/music-box-game/backend/internal/controllers"
	"github.com/domnikl/music-box-game/backend/internal/middlewares"
	"github.com/domnikl/music-box-game/backend/internal/repositories"
	"github.com/domnikl/music-box-game/backend/internal/services"
	"github.com/domnikl/music-box-game/backend/internal/spotify"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var token string

func setupSpotify(e *echo.Echo, db *gorm.DB, spotify *spotify.Spotify) {
	userRepository := repositories.NewUserRepository(db)
	apiUserAuthMiddleware := middlewares.NewAPIUserAuthMiddleware(userRepository)
	spotifyMiddleware := middlewares.NewSpotifyMiddleware()

	userService := services.NewUserService(userRepository)
	controller := controllers.NewSpotifyController(userService, spotify)

	g := e.Group("/spotify")
	g.Use(apiUserAuthMiddleware.IsAuthenticated)
	g.GET("/auth", controller.Auth)
	g.GET("/callback", controller.Callback)

	needsSpotifyToken := g.Group("")
	needsSpotifyToken.Use(spotifyMiddleware.HasToken)
	needsSpotifyToken.GET("/playlists", controller.GetPlaylists)
	needsSpotifyToken.GET("/playlists/:id", controller.GetPlaylist)
	needsSpotifyToken.GET("/devices", controller.GetDevices)
	needsSpotifyToken.GET("/currently-playing", controller.GetCurrentlyPlaying)
	needsSpotifyToken.POST("/player/next", controller.Next)
	needsSpotifyToken.POST("/player/pause", controller.Pause)
	needsSpotifyToken.POST("/player/play", controller.Play)
}
