package validate

import (
	"fmt"
	"regexp"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

var phonePattern = regexp.MustCompile(`^[0-9+ -]{6,20}$`)

func FarmerRecord(record domain.FarmerRecord) []string {
	errors := make([]string, 0, 6)
	if strings.TrimSpace(record.HouseholdHead) == "" {
		errors = append(errors, "household head is required")
	}
	if strings.TrimSpace(record.VillageGroup) == "" {
		errors = append(errors, "village group is required")
	}
	if record.CultivatedArea <= 0 {
		errors = append(errors, "cultivated area must be positive")
	}
	if strings.TrimSpace(record.MainCrop) == "" {
		errors = append(errors, "main crop is required")
	}
	if !phonePattern.MatchString(strings.TrimSpace(record.Phone)) {
		errors = append(errors, "phone format is invalid")
	}
	if !Date(record.LastVisit) {
		errors = append(errors, "last visit must use YYYY-MM-DD")
	}
	return errors
}

func VisitReport(report domain.VisitReport) []string {
	errors := make([]string, 0, 4)
	if strings.TrimSpace(report.FarmerID) == "" {
		errors = append(errors, "farmer id is required")
	}
	if !Date(report.VisitDate) {
		errors = append(errors, "visit date must use YYYY-MM-DD")
	}
	if strings.TrimSpace(report.Officer) == "" {
		errors = append(errors, "officer is required")
	}
	if strings.TrimSpace(report.Findings) == "" {
		errors = append(errors, "findings are required")
	}
	return errors
}

func ReviewCase(review domain.ReviewCase) []string {
	errors := make([]string, 0, 3)
	if strings.TrimSpace(review.FarmerID) == "" {
		errors = append(errors, "farmer id is required")
	}
	if strings.TrimSpace(review.ReportID) == "" {
		errors = append(errors, "report id is required")
	}
	if strings.TrimSpace(review.Reviewer) == "" {
		errors = append(errors, "reviewer is required")
	}
	if review.Decision != domain.ReviewApproved && review.Decision != domain.ReviewRejected {
		errors = append(errors, "decision must be approved or rejected")
	}
	return errors
}

func CollaborationNote(note domain.CollaborationNote) []string {
	errors := make([]string, 0, 3)
	if strings.TrimSpace(note.FarmerID) == "" {
		errors = append(errors, "farmer id is required")
	}
	if strings.TrimSpace(note.Author) == "" {
		errors = append(errors, "author is required")
	}
	if strings.TrimSpace(note.Body) == "" {
		errors = append(errors, "body is required")
	}
	if note.Visibility != "team" && note.Visibility != "private" {
		errors = append(errors, "visibility must be team or private")
	}
	return errors
}

func ArchiveEntry(entry domain.ArchiveEntry) []string {
	errors := make([]string, 0, 3)
	if strings.TrimSpace(entry.FarmerID) == "" {
		errors = append(errors, "farmer id is required")
	}
	if strings.TrimSpace(entry.ArchivedBy) == "" {
		errors = append(errors, "archived by is required")
	}
	if !Date(entry.ArchivedOn) {
		errors = append(errors, "archive date must use YYYY-MM-DD")
	}
	if strings.TrimSpace(entry.Reason) == "" {
		errors = append(errors, "archive reason is required")
	}
	return errors
}

func Date(value string) bool {
	if len(value) != 10 || value[4] != '-' || value[7] != '-' {
		return false
	}
	for index, value := range value {
		if index == 4 || index == 7 {
			continue
		}
		if value < '0' || value > '9' {
			return false
		}
	}
	return value[:4] != "0000" && value[5:7] != "00" && value[8:] != "00"
}

func RequireUniqueIDs(records []domain.FarmerRecord) error {
	seen := make(map[string]struct{}, len(records))
	for index, record := range records {
		if record.ID == "" {
			continue
		}
		if _, exists := seen[record.ID]; exists {
			return fmt.Errorf("duplicate farmer id at row %d: %s", index+1, record.ID)
		}
		seen[record.ID] = struct{}{}
	}
	return nil
}

func HasErrors(errors []string) bool {
	return len(errors) > 0
}
