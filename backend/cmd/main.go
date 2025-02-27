package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/domnikl/music-box-game/backend/internal/repositories"
	"github.com/domnikl/music-box-game/backend/internal/routes"
	"github.com/domnikl/music-box-game/backend/internal/services"
	"github.com/domnikl/music-box-game/backend/internal/spotify"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetRequiredEnvValue(name string) string {
	value := os.Getenv(name)
	if value == "" {
		slog.Error(fmt.Sprintf("Environment variable %s is not set", name))
		os.Exit(1)
	}

	return value
}

func main() {
	spotifyClientID := GetRequiredEnvValue("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := GetRequiredEnvValue("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURL := GetRequiredEnvValue("SPOTIFY_REDIRECT_URI")

	db, err := gorm.Open(postgres.Open(GetRequiredEnvValue("DB_DSN")), &gorm.Config{})
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
		services.NewUserService(repositories.NewUserRepository(db)),
	)

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))
	routes.Setup(e, db, spotifyClient)

	e.Logger.Fatal(e.Start(":8080"))
}
