package goterrauser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// User represents a user
type User struct {
	Logged    bool   `json:"logged"`
	UID       string `json:"uid"`
	APIKey    string `json:"apikey"`
	Password  string `json:"password"`
	Admin     bool   `json:"admin"`
	SuperUser bool   `json:"super"`
	Email     string `json:"email"`
	Kind      string `json:"kind"`                   // Kind of auth used by user (local, google, etc.)
	SSHPubKey string `json:"pub_key" bson:"pub_key"` // ssh public key
}

// AuthData is result struct for authentication with user data and an authentication token
type AuthData struct {
	User  User
	Token string
}

// Check checks X-API-Key authorization content and returns user info
func Check(apiKey string) (AuthData, error) {
	autData := AuthData{}

	url := os.Getenv("GOT_PROXY")
	if os.Getenv("GOT_PROXY_AUTH") != "" {
		url = os.Getenv("GOT_PROXY_AUTH")
	}

	client := &http.Client{}
	remote := []string{url, "auth", "api"}

	req, _ := http.NewRequest("GET", strings.Join(remote, "/"), nil)
	req.Header.Add("X-API-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return autData, errors.New("failed to contact auth service")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return autData, fmt.Errorf("auth error %d", resp.StatusCode)
	}
	respData := &AuthData{}
	json.NewDecoder(resp.Body).Decode(respData)
	return *respData, err
}
