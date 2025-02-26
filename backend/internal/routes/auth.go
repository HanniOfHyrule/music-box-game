package routes

import (
	"github.com/domnikl/music-box-game/backend/internal/controllers"
	"github.com/domnikl/music-box-game/backend/internal/repositories"
	"github.com/domnikl/music-box-game/backend/internal/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func setupAuth(e *echo.Echo, db *gorm.DB) {
	userService := services.NewUserService(repositories.NewUserRepository(db))
	authController := controllers.NewAuthController(userService)

	// TODO: secure this endpoint with a hardcoded secret
	e.POST("/auth", authController.InitAuth)
}
