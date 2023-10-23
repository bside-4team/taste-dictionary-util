package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to loading .env file")
	}

	conf := makeConfig()
	credBytes, err := base64.StdEncoding.DecodeString(conf.googleAPIKey)
	if err != nil {
		log.Fatalf("failed to decode google service account key")
	}

	gConfig, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("err: %v", err)
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

	index := -1
	for _, row := range sheet.Rows {
		// keyword := fmt.Sprintf("%s %s %s", row[0].Value, row[1].Value, row[2].Value)
		keyword := fmt.Sprintf("%s %s", row[3].Value, row[2].Value)
		index++
		cusines, err := searchCusineByKeyWord(conf, keyword)
		if err != nil {
			log.Printf("keyword: %v, something goes wrong with kakao search api\n", keyword)
			continue
		}

		if len(cusines) == 0 {
			log.Printf("keyword: %v, no result\n", keyword)
			continue
		}

		cusine := cusines[0]
		cusineFieldMap, ok := cusine.(map[string]interface{})
		if !ok {
			fmt.Println("Error: Invalid document format")
			continue
		}

		addressName := cusineFieldMap["address_name"].(string)
		categoryName := cusineFieldMap["category_name"].(string)
		id := cusineFieldMap["id"].(string)
		phone := cusineFieldMap["phone"].(string)
		placeName := cusineFieldMap["place_name"].(string)
		placeUrl := cusineFieldMap["place_url"].(string)
		roadAddressName := cusineFieldMap["road_address_name"].(string)
		x := cusineFieldMap["x"].(string)
		y := cusineFieldMap["y"].(string)

		sheet.Update(index, 8, placeName)
		sheet.Update(index, 9, addressName)
		sheet.Update(index, 10, id)
		sheet.Update(index, 11, phone)
		sheet.Update(index, 12, categoryName)
		sheet.Update(index, 13, placeUrl)
		sheet.Update(index, 14, roadAddressName)
		sheet.Update(index, 15, x)
		sheet.Update(index, 16, y)
	}
	err = sheet.Synchronize()
	if err != nil {
		log.Printf("keyword (%v): failed to sync sheet\n", err)
	}
}

func searchCusineByKeyWord(conf *config, keyword string) ([]interface{}, error) {
	encodedKeyWord := url.QueryEscape(keyword)
	url := fmt.Sprintf("https://dapi.kakao.com/v2/local/search/keyword.json?query=%s", encodedKeyWord)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request")
	}

	kakaoAPIKey := fmt.Sprintf("KakaoAK %s", conf.kakaoAPIKey)

	req.Header.Set("Authorization", kakaoAPIKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to making request")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to reading response body")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decoding JSON")
	}

	cusines := result["documents"].([]interface{})
	return cusines, nil
}
