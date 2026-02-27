package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

var GoogleOauthConfig *oauth2.Config

func InitGoogle() error {
	// Load credentials.json (from Google Cloud Console)
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("unable to read credentials.json: %v", err)
	}

	// Parse config with Drive scope
	conf, err := google.ConfigFromJSON(b,
		drive.DriveScope,
		"https://www.googleapis.com/auth/userinfo.email",
	)
	if err != nil {
		return fmt.Errorf("unable to parse credentials.json: %v", err)
	}

	GoogleOauthConfig = conf
	return nil
}

// Get Drive service with token.json
func GetDriveService() (*drive.Service, error) {
	ctx := context.Background()

	// Read token.json (refresh + access token)
	f, err := os.Open("token.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open token.json: %v", err)
	}
	defer f.Close()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, fmt.Errorf("unable to decode token.json: %v", err)
	}

	client := GoogleOauthConfig.Client(ctx, tok)
	return drive.New(client)
}
