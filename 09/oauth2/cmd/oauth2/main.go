package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	clientID      = "139d9832c5ded86ef1f8"
	clientSecret  = ""
	code          = ""
	tokenEndpoint = "https://github.com/login/oauth/access_token"
	userEndpoint  = "https://api.github.com/user"
)

func prepareRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	query := req.URL.Query()
	query.Set("client_id", clientID)
	query.Set("client_secret", clientSecret)
	query.Set("code", code)
	req.URL.RawQuery = query.Encode()

	return req, nil
}

type GitHubResponse struct {
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func getAccessToken() (string, error) {
	req, err := prepareRequest()
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return "", fmt.Errorf("expected status %d, got %d", http.StatusOK, code)
	}

	var ghResp GitHubResponse
	if err = json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return "", err
	}

	if errMsg := ghResp.Error; errMsg != "" {
		return "", fmt.Errorf("%s: %s", errMsg, ghResp.ErrorDescription)
	}

	return ghResp.AccessToken, nil
}

type User struct {
	ID              int       `json:"id"`
	Login           string    `json:"login"`
	Name            string    `json:"name"`
	Company         string    `json:"company"`
	AvatarURL       string    `json:"avatar_url"`
	HTMLUrl         string    `json:"html_url"`
	TwitterUserName string    `json:"twitter_username"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func getUser(accessToken string) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, userEndpoint, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("expected status %d, got %d", http.StatusOK, code)
	}

	var user User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	accessToken, err := getAccessToken()
	if err != nil {
		l.Fatal("http request failed", zap.Error(err))
	}
	l.Debug("access token", zap.String("token", accessToken))

	user, err := getUser(accessToken)
	if err != nil {
		l.Fatal("getting user", zap.Error(err))
	}

	l.Info("user info", zap.Any("user_details", user))
}
