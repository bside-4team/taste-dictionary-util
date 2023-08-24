package gsheet

import (
	"context"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

type GsheetService struct {
	gsheetService *sheets.Service
	sheetId       string
}

func NewGSheetSerive(gsheetService *sheets.Service, sheetId string) (*GsheetService, error) {
	ctx := context.Background()
	sheetSerivce, err := sheets.NewService(ctx)

	if err != nil {
		return nil, fmt.Errorf("error creating new GsheetService: %v", err)
	}

	return &GsheetService{
		gsheetService: sheetSerivce,
		sheetId:       sheetId,
	}, nil
}

func (g *GsheetService) ReadSheetData(sheetName string) ([][]interface{}, error) {
	res, err := g.gsheetService.Spreadsheets.Values.Get(g.sheetId, sheetName).Do()
	if err != nil {
		return nil, err
	}

	return res.Values, nil
}

func (g *GsheetService) UpdateSheetData(sheetName string, data [][]interface{}) error {
	_, err := g.gsheetService.Spreadsheets.Values.Update(g.sheetId, sheetName, &sheets.ValueRange{
		Values: data,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}
