package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	foodDataFile = "menu.csv"
	outputFile   = "food.json"
)

type Food struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Category []string `json:"category"`
	Keyword  []string `json:"keyword"`
}

type Data struct {
	Data     []Food `json:"data"`
	Metadata struct {
		Total int `json:"total"`
	} `json:"meta"`
}

func main() {
	records, err := readCsv(foodDataFile)

	if err != nil {
		log.Fatal("failed to read csv file")
	}

	jsonFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer jsonFile.Close()

	if err = ParseToJson(records, jsonFile); err != nil {
		fmt.Println("failed to parse json")
	}

}

func readCsv(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv file")
	}

	return records, nil
}

func ParseToJson(records [][]string, jsonFile io.Writer) error {
	foodItems := make([]Food, 0)

	for index, row := range records {
		name := row[0]
		categories := split(row[1], ", ")
		keywords := split(row[2], ", ")

		food := &Food{
			Id:       index,
			Name:     name,
			Category: categories,
			Keyword:  keywords,
		}

		foodItems = append(foodItems, *food)
	}

	Data := Data{
		Data: foodItems,
		Metadata: struct {
			Total int `json:"total"`
		}{
			Total: len(foodItems),
		},
	}

	// data, err := json.MarshalIndent(Data, "", " ")
	// if err != nil {
	// 	return fmt.Errorf("failed to Marshal")
	// }
	// fmt.Println(string(data))

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(Data); err != nil {
		return fmt.Errorf("error writing JSON file")
	}

	return nil
}

func split(str, sep string) []string {
	if len(str) == 0 {
		return []string{}
	}
	return strings.Split(str, sep)
}
