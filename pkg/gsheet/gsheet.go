package gsheet

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

type config struct {
	kakaoAPIKey  string
	googleAPIKey string
}

func makeConfig() *config {
	return &config{
		kakaoAPIKey:  os.Getenv("KAKAO_API_KEY"),
		googleAPIKey: os.Getenv("GOOGLE_API_KEY"),
	}
}

var (
	sheetId = "1PfboVci0tyuw-JdoL6v1hoCVL2eoim2nkeZDm5fjX3Y"
)

func GetSheet(sheetid int) (*spreadsheet.Sheet, error) {
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode google service account key")
	}

	gConfig, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	client := gConfig.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(sheetId)

	if err != nil {
		log.Fatalf("failed to fetch spreadsheet")
	}

	sheet, err := spreadsheet.SheetByIndex(0)
	if err != nil {
		log.Fatal("failed to init google sheet service")
	}

	return sheet, nil
}
