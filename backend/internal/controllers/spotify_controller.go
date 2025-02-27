package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/domnikl/music-box-game/backend/internal/crypto"
	"github.com/domnikl/music-box-game/backend/internal/models"
	"github.com/domnikl/music-box-game/backend/internal/services"
	"github.com/domnikl/music-box-game/backend/internal/spotify"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type SpotifyController struct {
	userService *services.UserService
	spotify     *spotify.Spotify
}

func NewSpotifyController(userService *services.UserService, spotify *spotify.Spotify) *SpotifyController {
	return &SpotifyController{userService: userService, spotify: spotify}
}

func (s *SpotifyController) Auth(c echo.Context) error {
	user := c.Get("user").(*models.User)
	state := fmt.Sprintf("%s:%s", crypto.RandomAlphaNumericString(32), user.APIToken)
	url := s.spotify.AuthURL(state)

	sess, err := session.Get("session", c)
	if err != nil {
		slog.Error("Failed to save session: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["state"] = state

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		slog.Error("Failed to save session: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.Redirect(http.StatusFound, url)
}

func (s *SpotifyController) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	sess, err := session.Get("session", c)
	if err != nil {
		slog.Error("Failed to get session: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	sessionState, ok := sess.Values["state"].(string)
	if !ok || state != sessionState {
		slog.Error("State mismatch: " + state + " != " + sessionState)
		return c.String(http.StatusBadRequest, "State mismatch")
	}

	// clear state
	sess.Values["state"] = nil
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		slog.Error("Failed to saved session: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	t, err := s.spotify.Exchange(code)
	if err != nil {
		slog.Error("Failed to exchange code for token: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	user := c.Get("user").(*models.User)
	user.SpotifyToken = t.AccessToken
	user.SpotifyRefreshToken = t.RefreshToken

	err = s.userService.UpdateUser(user)
	if err != nil {
		slog.Error("Failed to update user: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.String(http.StatusOK, "You can now close this window")
}

func (s *SpotifyController) GetPlaylist(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("user").(*models.User)

	playlist, err := s.spotify.GetPlaylist(user, id)

	if err != nil {
		slog.Error("Failed to get playlist: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, playlist)
}

func (s *SpotifyController) GetPlaylists(c echo.Context) error {
	queryLimit := c.QueryParam("limit")
	if queryLimit == "" {
		queryLimit = "50"
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		slog.Error("Failed to parse limit: " + err.Error())
		return c.String(http.StatusBadRequest, "Invalid limit")
	}

	queryOffset := c.QueryParam("offset")
	if queryOffset == "" {
		queryOffset = "0"
	}

	offset, err := strconv.Atoi(queryOffset)
	if err != nil {
		slog.Error("Failed to parse offset: " + err.Error())
		return c.String(http.StatusBadRequest, "Invalid offset")
	}

	user := c.Get("user").(*models.User)
	playlists, err := s.spotify.GetPlaylists(user, limit, offset)

	if err != nil {
		slog.Error("Failed to get playlists: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, playlists)
}

func (s *SpotifyController) Next(c echo.Context) error {
	user := c.Get("user").(*models.User)
	err := s.spotify.Next(user)

	if err != nil {
		slog.Error("Failed to skip track: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *SpotifyController) Pause(c echo.Context) error {
	user := c.Get("user").(*models.User)
	err := s.spotify.Pause(user)

	if err != nil {
		slog.Error("Failed to pause playback: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *SpotifyController) Play(c echo.Context) error {
	type playRequest struct {
		DeviceID   string `json:"device_id,omitempty"`
		PlaylistID string `json:"playlist_id,omitempty"`
	}

	user := c.Get("user").(*models.User)

	var req playRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request: " + err.Error())
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	err := s.spotify.Play(user, req.DeviceID, req.PlaylistID)

	if err != nil {
		slog.Error("Failed to start playback: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *SpotifyController) GetDevices(c echo.Context) error {
	user := c.Get("user").(*models.User)

	devices, err := s.spotify.GetDevices(user)

	if err != nil {
		slog.Error("Failed to get devices: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, devices)
}

func (s *SpotifyController) GetCurrentlyPlaying(c echo.Context) error {
	user := c.Get("user").(*models.User)

	currentlyPlaying, err := s.spotify.GetCurrentlyPlaying(user)

	if err != nil {
		slog.Error("Failed to get currently playing: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, currentlyPlaying)
}
