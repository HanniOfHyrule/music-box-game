package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/domnikl/music-box-game/backend/internal/repositories"
	"github.com/labstack/echo/v4"
)

type APIUserAuthMiddleware struct {
	userRepository *repositories.UserRepository
}

func NewAPIUserAuthMiddleware(userRepository *repositories.UserRepository) APIUserAuthMiddleware {
	return APIUserAuthMiddleware{userRepository: userRepository}
}

func (m APIUserAuthMiddleware) IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiToken, err := apiTokenFromHeader(c)
		if err != nil {
			return err
		}

		if apiToken == "" {
			apiToken = apiTokenFromQuery(c)
		}

		if apiToken == "" {
			apiToken = apiTokenFromState(c)
		}

		if apiToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Unauthorized",
			})
		}

		user, err := m.userRepository.FindUserByAPIToken(apiToken)
		if err != nil || user == nil {
			slog.Error("Failed to find user by API token", "error", err)

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Unauthorized",
			})
		}

		c.Set("user", user)

		return next(c)
	}
}

func apiTokenFromHeader(c echo.Context) (string, error) {
	apiToken := c.Request().Header.Get("Authorization")

	if apiToken != "" && !strings.HasPrefix(apiToken, "Bearer ") {
		return "", c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}

	return strings.TrimPrefix(apiToken, "Bearer "), nil
}

func apiTokenFromQuery(c echo.Context) string {
	return c.QueryParam("api_token")
}

func apiTokenFromState(c echo.Context) string {
	state := c.QueryParam("state")
	if state == "" {
		return ""
	}

	parts := strings.Split(state, ":")
	if len(parts) != 2 {
		return ""
	}

	return strings.Split(c.QueryParam("state"), ":")[1]
}
