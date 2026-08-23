package quality

import (
	"fmt"
	"sort"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/validate"
)

type Finding struct {
	Code     string
	Field    string
	Severity string
	Message  string
}

type Assessment struct {
	RecordID string
	Findings []Finding
	Score    int
}

func Assess(record domain.FarmerRecord) Assessment {
	findings := make([]Finding, 0)
	if strings.TrimSpace(record.HouseholdHead) == "" {
		findings = append(findings, Finding{Code: "HEAD_EMPTY", Field: "household_head", Severity: "error", Message: "household head is empty"})
	}
	if strings.TrimSpace(record.VillageGroup) == "" {
		findings = append(findings, Finding{Code: "VILLAGE_EMPTY", Field: "village_group", Severity: "error", Message: "village group is empty"})
	}
	if record.CultivatedArea > 50 {
		findings = append(findings, Finding{Code: "AREA_REVIEW", Field: "cultivated_area", Severity: "warning", Message: "area needs manual review"})
	}
	if strings.TrimSpace(record.MainCrop) == "" {
		findings = append(findings, Finding{Code: "CROP_EMPTY", Field: "main_crop", Severity: "error", Message: "main crop is empty"})
	}
	if strings.TrimSpace(record.Phone) == "" {
		findings = append(findings, Finding{Code: "PHONE_EMPTY", Field: "phone", Severity: "warning", Message: "phone is empty"})
	}
	if !validate.Date(record.LastVisit) {
		findings = append(findings, Finding{Code: "VISIT_INVALID", Field: "last_visit", Severity: "error", Message: "last visit is invalid"})
	}
	score := qualityScore(findings)
	return Assessment{RecordID: record.ID, Findings: findings, Score: score}
}

func qualityScore(findings []Finding) int {
	score := 100
	for _, finding := range findings {
		switch finding.Severity {
		case "error":
			score -= 25
		case "warning":
			score -= 10
		default:
			score -= 1
		}
	}
	if score < 0 {
		return 0
	}
	return score
}

func AssessAll(records []domain.FarmerRecord) []Assessment {
	result := make([]Assessment, 0, len(records))
	for _, record := range records {
		result = append(result, Assess(record))
	}
	sort.SliceStable(result, func(left, right int) bool {
		if result[left].Score == result[right].Score {
			return result[left].RecordID < result[right].RecordID
		}
		return result[left].Score < result[right].Score
	})
	return result
}

func Errors(assessment Assessment) []Finding {
	result := make([]Finding, 0)
	for _, finding := range assessment.Findings {
		if finding.Severity == "error" {
			result = append(result, finding)
		}
	}
	return result
}

func Warnings(assessment Assessment) []Finding {
	result := make([]Finding, 0)
	for _, finding := range assessment.Findings {
		if finding.Severity == "warning" {
			result = append(result, finding)
		}
	}
	return result
}

func CanPublish(assessment Assessment) bool {
	return len(Errors(assessment)) == 0
}

func FindingText(finding Finding) string {
	return fmt.Sprintf("%s:%s:%s", finding.Code, finding.Field, finding.Message)
}

func SummaryText(assessment Assessment) string {
	return fmt.Sprintf("%s score=%d findings=%d", assessment.RecordID, assessment.Score, len(assessment.Findings))
}

func FieldsWithFindings(assessment Assessment) []string {
	seen := make(map[string]struct{})
	fields := make([]string, 0)
	for _, finding := range assessment.Findings {
		if _, exists := seen[finding.Field]; exists {
			continue
		}
		seen[finding.Field] = struct{}{}
		fields = append(fields, finding.Field)
	}
	sort.Strings(fields)
	return fields
}
