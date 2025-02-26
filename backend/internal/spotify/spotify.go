package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

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

func (s *Spotify) Playlists(accessToken string, limit, offset int) (*PlaylistsResponse, error) {
	// TODO: follow pagination

	httpClient := &http.Client{}
	url, err := url.Parse("https://api.spotify.com/v1/me/playlists?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			"Authorization": []string{"Bearer " + accessToken},
		},
	}

	resp, err := httpClient.Do(request)
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

func (s *Spotify) Playlist(accessToken, id string) (*PlaylistResponse, error) {
	httpClient := &http.Client{}
	url, err := url.Parse("https://api.spotify.com/v1/playlists/" + id)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			"Authorization": []string{"Bearer " + accessToken},
		},
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get playlist: %s %s", resp.Status, string(body))
	}

	// json decode
	var playlist PlaylistResponse
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

func (s *Spotify) Next(accessToken string) error {
	httpClient := &http.Client{}
	url, err := url.Parse("https://api.spotify.com/v1/me/player/next")
	if err != nil {
		return err
	}

	request := &http.Request{
		Method: "POST",
		URL:    url,
		Header: http.Header{
			"Authorization": []string{"Bearer " + accessToken},
		},
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to skip track: %s %s", resp.Status, string(body))
	}

	return nil
}
