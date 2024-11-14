package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/JonathanDonald03/spotify-tui/spotifyauth"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Some Error Occured: %s \n", err)
	}

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			// Add other scopes as needed
		),
	)

	if err != nil {
		fmt.Printf("Error fetching playlists: %s \n", err)
	}

}

func get_client_playlists(access_token string) (*PlaylistResponse, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header with the access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response JSON into a PlaylistsResponse struct
	var playlistsResponse PlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlistsResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &playlistsResponse, nil
}

type Playlist struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	PublicAccess bool   `json:"public"`
	TracksTotal  int    `json:"tracks_total"`
}

type PlaylistResponse struct {
	Items []Playlist `json:"items"`
}

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func get_profile(access_token string) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
		return
	}
	fmt.Println("Response Body:", string(body))
}
