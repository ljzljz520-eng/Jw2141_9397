package importer

import (
	"fmt"
	"io"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/service"
	"example.com/xiangzhenfarm/internal/validate"
)

type BatchFile struct {
	Name  string
	Input io.Reader
}

type BatchResult struct {
	Files    int
	Imported int
	Rejected int
	Rows     []domain.ImportRow
}

func ProcessFiles(files []BatchFile, registry *service.Registry) (BatchResult, error) {
	if len(files) == 0 {
		return BatchResult{}, fmt.Errorf("at least one import file is required")
	}
	result := BatchResult{Rows: make([]domain.ImportRow, 0)}
	importer := NewCSVImporter()
	for _, file := range files {
		summary, err := importer.Read(file.Input)
		if err != nil {
			return result, fmt.Errorf("process %s: %w", file.Name, err)
		}
		result.Files++
		result.Imported += summary.Accepted
		result.Rejected += summary.Rejected
		result.Rows = append(result.Rows, summary.Rows...)
		for _, record := range importer.ValidRows(summary) {
			if _, err := registry.CreateFarmer(record); err != nil {
				return result, fmt.Errorf("persist %s: %w", file.Name, err)
			}
		}
	}
	if err := validate.RequireUniqueIDs(uniqueRecords(result.Rows)); err != nil {
		return result, err
	}
	return result, nil
}

func uniqueRecords(rows []domain.ImportRow) []domain.FarmerRecord {
	result := make([]domain.FarmerRecord, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		if len(row.Errors) > 0 {
			continue
		}
		if _, exists := seen[row.Record.ID]; exists {
			continue
		}
		seen[row.Record.ID] = struct{}{}
		result = append(result, row.Record)
	}
	return result
}

func CountFiles(result BatchResult) int {
	return result.Files
}

func AcceptedRows(result BatchResult) int {
	return result.Imported
}

func RejectedRows(result BatchResult) int {
	return result.Rejected
}
