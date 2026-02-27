package utils

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// GetSheetsService initializes Google Sheets API client
func GetSheetsService() *sheets.Service {
	ctx := context.Background()

	// Load credentials from JSON key file (OAuth2 service account)
	creds, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Sheets client: %v", err)
	}

	return srv
}
