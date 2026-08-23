package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/validate"
)

var requiredHeaders = []string{"id", "household_head", "village_group", "cultivated_area", "main_crop", "phone", "last_visit"}

type CSVImporter struct{}

func NewCSVImporter() *CSVImporter {
	return &CSVImporter{}
}

func (i *CSVImporter) Read(input io.Reader) (domain.ImportSummary, error) {
	reader := csv.NewReader(input)
	header, err := reader.Read()
	if err != nil {
		return domain.ImportSummary{}, fmt.Errorf("read csv header: %w", err)
	}
	if err := validateHeaders(header); err != nil {
		return domain.ImportSummary{}, err
	}
	summary := domain.ImportSummary{Rows: make([]domain.ImportRow, 0)}
	for line := 2; ; line++ {
		row, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return summary, fmt.Errorf("read csv row %d: %w", line, readErr)
		}
		record, parseErr := parseRow(row, line)
		entry := domain.ImportRow{Line: line, Record: record}
		if parseErr != nil {
			entry.Errors = []string{parseErr.Error()}
			summary.Rejected++
		} else if validation := validate.FarmerRecord(record); validate.HasErrors(validation) {
			entry.Errors = validation
			summary.Rejected++
		} else {
			summary.Accepted++
		}
		summary.Rows = append(summary.Rows, entry)
	}
	return summary, nil
}

func validateHeaders(header []string) error {
	if len(header) != len(requiredHeaders) {
		return fmt.Errorf("csv requires %d columns", len(requiredHeaders))
	}
	for index, expected := range requiredHeaders {
		if strings.TrimSpace(strings.ToLower(header[index])) != expected {
			return fmt.Errorf("csv column %d must be %s", index+1, expected)
		}
	}
	return nil
}

func parseRow(row []string, line int) (domain.FarmerRecord, error) {
	if len(row) != len(requiredHeaders) {
		return domain.FarmerRecord{}, fmt.Errorf("line %d has %d columns", line, len(row))
	}
	area, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil {
		return domain.FarmerRecord{}, fmt.Errorf("line %d area: %w", line, err)
	}
	return domain.FarmerRecord{
		ID:             strings.TrimSpace(row[0]),
		HouseholdHead:  strings.TrimSpace(row[1]),
		VillageGroup:   strings.TrimSpace(row[2]),
		CultivatedArea: area,
		MainCrop:       strings.TrimSpace(row[4]),
		Phone:          strings.TrimSpace(row[5]),
		LastVisit:      strings.TrimSpace(row[6]),
		Status:         domain.RecordActive,
		CreatedFrom:    fmt.Sprintf("csv-line-%d", line),
	}, nil
}

func (i *CSVImporter) ValidRows(summary domain.ImportSummary) []domain.FarmerRecord {
	result := make([]domain.FarmerRecord, 0, summary.Accepted)
	for _, row := range summary.Rows {
		if len(row.Errors) == 0 {
			result = append(result, row.Record)
		}
	}
	return result
}

func (i *CSVImporter) ErrorRows(summary domain.ImportSummary) []domain.ImportRow {
	result := make([]domain.ImportRow, 0, summary.Rejected)
	for _, row := range summary.Rows {
		if len(row.Errors) > 0 {
			result = append(result, row)
		}
	}
	return result
}

func (i *CSVImporter) Header() []string {
	return append([]string(nil), requiredHeaders...)
}
