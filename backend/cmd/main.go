package main

import (
	"log/slog"
	"os"

	"github.com/domnikl/music-box-game/backend/internal/routes"
	"github.com/domnikl/music-box-game/backend/internal/spotify"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURL := os.Getenv("SPOTIFY_REDIRECT_URI")

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// run migrations
	migrationDir := os.Getenv("GOOSE_MIGRATION_DIR")

	database, err := db.DB()
	if err != nil {
		slog.Error("Failed to connect to database: " + err.Error())
		os.Exit(1)
	}

	if err = goose.Up(database, migrationDir); err != nil {
		slog.Error("Failed to run migrations: " + err.Error())
		os.Exit(1)
	}

	spotifyClient := spotify.NewSpotify(
		spotifyClientID,
		spotifyClientSecret,
		spotifyRedirectURL,
	)

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))
	routes.Setup(e, db, spotifyClient)

	e.Logger.Fatal(e.Start(":8080"))
}
