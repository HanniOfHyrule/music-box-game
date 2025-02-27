package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/domnikl/music-box-game/backend/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type Spotify struct {
	oauthConfig oauth2.Config
}

func NewSpotify(clientID, clientSecret, redirectURL string) *Spotify {
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user-read-email", "playlist-read-private", "playlist-read-collaborative", "user-modify-playback-state"},
		Endpoint:     spotify.Endpoint,
		RedirectURL:  redirectURL,
	}

	return &Spotify{oauthConfig: conf}
}

func (s *Spotify) AuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Spotify) Exchange(code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(context.Background(), code)
}

type PlaylistsResponse struct {
	Items []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Public      bool   `json:"public"`
		Images      []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Tracks struct {
			Total int `json:"total"`
		} `json:"tracks"`
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"items"`
}

func (s *Spotify) Playlists(user *models.User, limit, offset int) (*PlaylistsResponse, error) {
	resp, err := s.doRequest(http.MethodGet, "/me/playlists?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), user)
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

func (s *Spotify) Playlist(user *models.User, id string) (*PlaylistResponse, error) {
	resp, err := s.doRequest(http.MethodGet, "/playlists/"+id, user)
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
	resp, err := s.doRequest(http.MethodPost, "/me/player/next", user)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to skip track: %s %s", resp.Status, string(body))
	}

	return nil
}

func (s *Spotify) doRequest(method string, path string, user *models.User) (*http.Response, error) {
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
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return resp, err
	}

	// TODO: handle token refresh automatically
	if resp.StatusCode == http.StatusUnauthorized {
		return resp, fmt.Errorf("unauthorized")
	}

	return resp, nil
}
