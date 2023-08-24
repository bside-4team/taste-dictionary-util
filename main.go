package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	kakaoAPIKey string
}

func makeConfig() *config {
	return &config{
		kakaoAPIKey: os.Getenv("KAKAO_API_KEY"),
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to loading .env file")
	}

	conf := makeConfig()

	// TODO: login google account
	// Init google sheet service
	// Iterate sheet data
	// concatenate 지역 (A# ~), 도시명 (B# ~) 포털검색명 (F# ~)
	// call searchCusineByKeyWord
	// update sheet data H# (latitude), #I (longitude)

	keyword := "강남구 구찌라꾸"
	cusines, err := searchCusineByKeyWord(conf, keyword)

	if err != nil {
		log.Fatalf("failed to search cusine by keyword")
		os.Exit(1)
	}

	fmt.Println(cusines)
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
