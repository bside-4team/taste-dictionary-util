package main

import "os"

type Config struct {
	kakaoAPIKey  string
	googleAPIKey string
}

func MakeConfig() *Config {
	return &Config{
		kakaoAPIKey:  os.Getenv("KAKAO_API_KEY"),
		googleAPIKey: os.Getenv("GOOGLE_API_KEY"),
	}
}
