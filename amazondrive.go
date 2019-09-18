package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/toqueteos/webbrowser"

	"golang.org/x/oauth2"
)

type ACDriveStore struct {
	Config     oauth2.Config `json:"oauth_config"`
	OAuthToken oauth2.Token  `json:"oauth_token"`
}

func (acd *ACDriveStore) Setup() bool {
	config := receiveConfig()

	token, err := getACDriveTokenFromWeb(config)
	if err != nil {
		color.Red("Unable to get client token: %v", err)
		return false
	}

	acd.OAuthToken = *token

	return false
}

func getACDriveTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("go-acd", oauth2.AccessTypeOffline)
	webbrowser.Open(authURL)

	color.Yellow("Enter Auth Code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		color.Red("Unable to read authentication code %v", err)
		return nil, err
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		color.Red("Unable to retrieve token from web %v", err)
		return tok, err
	}

	return tok, nil
}

func receiveConfig() *oauth2.Config {
	var appKey, appSecret string

	appKey, appSecret = AmazonDriveClientKey, AmazonDriveClientSecret

	return &oauth2.Config{
		ClientID:     appKey,
		ClientSecret: appSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.amazon.com/ap/oa",
			TokenURL: "https://api.amazon.com/auth/o2/token",
		},
		RedirectURL: "https://go-acd.appspot.com/oauth",
		Scopes:      []string{"clouddrive:read", "clouddrive:write"},
	}
}
