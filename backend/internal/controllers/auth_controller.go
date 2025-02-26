package controllers

import (
	"github.com/domnikl/music-box-game/backend/internal/crypto"
	"github.com/domnikl/music-box-game/backend/internal/models"
	"github.com/domnikl/music-box-game/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{userService: userService}
}

func (a *AuthController) InitAuth(c echo.Context) error {
	// creates a new user with a random secret
	user := &models.User{APIToken: crypto.RandomAlphaNumericString(64)}

	err := a.userService.CreateUser(user)
	if err != nil {
		return c.JSON(500, map[string]string{
			"message": "could not create user",
		})
	}

	return c.JSON(200, user)
}
