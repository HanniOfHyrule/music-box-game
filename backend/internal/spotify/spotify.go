package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/domnikl/music-box-game/backend/internal/models"
	"github.com/domnikl/music-box-game/backend/internal/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type Spotify struct {
	oauthConfig oauth2.Config
	userService *services.UserService
}

func NewSpotify(clientID, clientSecret, redirectURL string, userService *services.UserService) *Spotify {
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user-read-email", "playlist-read-private", "playlist-read-collaborative", "user-modify-playback-state", "user-read-playback-state"},
		Endpoint:     spotify.Endpoint,
		RedirectURL:  redirectURL,
	}

	return &Spotify{oauthConfig: conf, userService: userService}
}

func (s *Spotify) AuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Spotify) Exchange(code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(context.Background(), code)
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Album struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	ReleaseDate string  `json:"release_date"`
	Images      []Image `json:"images"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type Item struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Public      bool    `json:"public"`
	Images      []Image `json:"images"`
	Tracks      struct {
		Total int `json:"total"`
	} `json:"tracks"`
	Type    string   `json:"type"`
	Album   Album    `json:"album"`
	Artists []Artist `json:"artists"`
}

type PlaylistsResponse struct {
	Items []Item `json:"items"`
}

func (s *Spotify) GetPlaylists(user *models.User, limit, offset int) (*PlaylistsResponse, error) {
	resp, err := s.doRequest(http.MethodGet, "/me/playlists?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), user, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get playlist: %s", resp.Status)
	}

	// json decode
	var playlists PlaylistsResponse
	err = json.NewDecoder(resp.Body).Decode(&playlists)
	if err != nil {
		return nil, err
	}

	return &playlists, nil
}

type PlaylistResponse struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Public        bool   `json:"public"`
	Collaborative bool   `json:"collaborative"`
	Images        []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Tracks struct {
		Total int `json:"total"`
		Items []struct {
			Track struct {
				Name  string `json:"name"`
				Album struct {
					Name        string `json:"name"`
					ReleaseDate string `json:"release_date"`
				} `json:"album"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (s *Spotify) GetPlaylist(user *models.User, id string) (*PlaylistResponse, error) {
	resp, err := s.doRequest(http.MethodGet, "/playlists/"+id, user, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get playlist: %s %s", resp.Status, string(body))
	}

	var playlist PlaylistResponse
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

func (s *Spotify) Next(user *models.User) error {
	resp, err := s.doRequest(http.MethodPost, "/me/player/next", user, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to skip track: %s %s", resp.Status, string(body))
	}

	return nil
}

func (s *Spotify) Pause(user *models.User) error {
	resp, err := s.doRequest(http.MethodPut, "/me/player/pause", user, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to pause playback: %s %s", resp.Status, string(body))
	}

	return nil
}

func (s *Spotify) Play(user *models.User, deviceID string, playlistID string) error {
	type playRequest struct {
		ContextURI string `json:"context_uri,omitempty"`
	}

	if playlistID != "" {
		playlistID = "spotify:playlist:" + playlistID
	}

	body, err := json.Marshal(playRequest{ContextURI: playlistID})
	if err != nil {
		return err
	}

	resp, err := s.doRequest(http.MethodPut, fmt.Sprintf("/me/player/play?device_id=%s", deviceID), user, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to start playback: %s %s", resp.Status, string(body))
	}

	return nil
}

type Device struct {
	ID             string `json:"id"`
	IsActive       bool   `json:"is_active"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	VolumePercent  int    `json:"volume_percent"`
	SupportsVolume bool   `json:"supports_volume"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

func (s *Spotify) GetDevices(user *models.User) ([]Device, error) {
	resp, err := s.doRequest(http.MethodGet, "/me/player/devices", user, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get devices: %s", resp.Status)
	}

	// json decode
	var response DevicesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Devices, nil
}

func (s *Spotify) doRequest(method string, path string, user *models.User, body []byte, refreshTokens ...bool) (*http.Response, error) {
	httpClient := &http.Client{}
	url, err := url.Parse("https://api.spotify.com/v1" + path)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: method,
		URL:    url,
		Header: http.Header{
			"Authorization": []string{"Bearer " + user.SpotifyToken},
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode == http.StatusUnauthorized && refreshTokens == nil {
		s.refreshToken(user)
		return s.doRequest(method, path, user, body, false)
	}

	return resp, nil
}

func (s *Spotify) refreshToken(user *models.User) error {
	tokenSource := s.oauthConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: user.SpotifyRefreshToken,
	})

	token, err := tokenSource.Token()
	if err != nil {
		return err
	}

	user.SpotifyToken = token.AccessToken
	user.SpotifyRefreshToken = token.RefreshToken

	// only update specific fields
	s.userService.UpdateUser(&models.User{
		ID:                  user.ID,
		SpotifyToken:        user.SpotifyToken,
		SpotifyRefreshToken: user.SpotifyRefreshToken,
	})

	return nil
}

type CurrentlyPlayingResponse struct {
	Device               *Device `json:"device"`
	RepeatState          string  `json:"repeat_state"`
	ShuffleState         bool    `json:"shuffle_state"`
	IsPlaying            bool    `json:"is_playing"`
	CurrentlyPlayingType string  `json:"currently_playing_type"`
	Timestamp            int64   `json:"timestamp"`
	ProgressMs           int     `json:"progress_ms"`
	Item                 Item    `json:"item"`
}

func (s *Spotify) GetCurrentlyPlaying(user *models.User) (*CurrentlyPlayingResponse, error) {
	resp, err := s.doRequest(http.MethodGet, "/me/player/currently-playing?market=DE", user, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get currently playing: %s", resp.Status)
	}

	// json decode
	var response CurrentlyPlayingResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
