package service

import (
	"fmt"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

func (r *Registry) ImportRows(rows []domain.ImportRow, actor string, date string) (int, error) {
	if len(rows) == 0 {
		return 0, fmt.Errorf("import contains no rows")
	}
	accepted := 0
	for _, row := range rows {
		if len(row.Errors) > 0 {
			continue
		}
		if _, err := r.CreateFarmer(row.Record); err != nil {
			return accepted, err
		}
		if _, err := r.RecordAudit("FarmerRecord", row.Record.ID, domain.AuditImported, actor, date, strings.Join(row.Errors, ";")); err != nil {
			return accepted, err
		}
		accepted++
	}
	if accepted == 0 {
		return 0, fmt.Errorf("import contains no valid farmer rows")
	}
	return accepted, nil
}

func ImportOutcome(accepted int, rejected int) string {
	if rejected == 0 {
		return "complete"
	}
	if accepted == 0 {
		return "rejected"
	}
	return "partial"
}

func ImportErrorSummary(rows []domain.ImportRow) string {
	parts := make([]string, 0)
	for _, row := range rows {
		if len(row.Errors) == 0 {
			continue
		}
		parts = append(parts, fmt.Sprintf("line %d: %s", row.Line, strings.Join(row.Errors, ", ")))
	}
	return strings.Join(parts, "; ")
}

func ImportLineCount(rows []domain.ImportRow) int {
	return len(rows)
}

func ImportRejectedCount(rows []domain.ImportRow) int {
	count := 0
	for _, row := range rows {
		if len(row.Errors) > 0 {
			count++
		}
	}
	return count
}

func ImportAcceptedCount(rows []domain.ImportRow) int {
	return len(rows) - ImportRejectedCount(rows)
}
