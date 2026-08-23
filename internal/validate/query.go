package validate

import (
	"sort"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

func NormalizeFilter(filter domain.SearchFilter) domain.SearchFilter {
	filter.HouseholdHead = strings.TrimSpace(strings.ToLower(filter.HouseholdHead))
	filter.VillageGroup = strings.TrimSpace(strings.ToLower(filter.VillageGroup))
	filter.MainCrop = strings.TrimSpace(strings.ToLower(filter.MainCrop))
	filter.Status = strings.TrimSpace(strings.ToLower(filter.Status))
	return filter
}

func MatchRecord(record domain.FarmerRecord, rawFilter domain.SearchFilter) bool {
	filter := NormalizeFilter(rawFilter)
	if filter.HouseholdHead != "" && !strings.Contains(strings.ToLower(record.HouseholdHead), filter.HouseholdHead) {
		return false
	}
	if filter.VillageGroup != "" && strings.ToLower(record.VillageGroup) != filter.VillageGroup {
		return false
	}
	if filter.MainCrop != "" && strings.ToLower(record.MainCrop) != filter.MainCrop {
		return false
	}
	if filter.Status != "" && strings.ToLower(domain.NormalizeStatus(record.Status, domain.RecordActive)) != filter.Status {
		return false
	}
	if filter.MinArea > 0 && record.CultivatedArea < filter.MinArea {
		return false
	}
	if filter.MaxArea > 0 && record.CultivatedArea > filter.MaxArea {
		return false
	}
	return true
}

func FilterRecords(records []domain.FarmerRecord, filter domain.SearchFilter) []domain.FarmerRecord {
	result := make([]domain.FarmerRecord, 0, len(records))
	for _, record := range records {
		if MatchRecord(record, filter) {
			result = append(result, record)
		}
	}
	return result
}

func SortRecords(records []domain.FarmerRecord, spec domain.SortSpec) []domain.FarmerRecord {
	result := append([]domain.FarmerRecord(nil), records...)
	sort.SliceStable(result, func(left, right int) bool {
		before := compareRecords(result[left], result[right], spec.Field)
		if spec.Reverse {
			return before > 0
		}
		return before < 0
	})
	return result
}

func compareRecords(left domain.FarmerRecord, right domain.FarmerRecord, field string) int {
	var l, r string
	switch strings.ToLower(field) {
	case "head", "household_head":
		l, r = left.HouseholdHead, right.HouseholdHead
	case "village", "village_group":
		l, r = left.VillageGroup, right.VillageGroup
	case "crop", "main_crop":
		l, r = left.MainCrop, right.MainCrop
	case "area", "cultivated_area":
		if left.CultivatedArea < right.CultivatedArea {
			return -1
		}
		if left.CultivatedArea > right.CultivatedArea {
			return 1
		}
		return 0
	default:
		l, r = left.LastVisit, right.LastVisit
	}
	if l < r {
		return -1
	}
	if l > r {
		return 1
	}
	return 0
}

func CropNames(records []domain.FarmerRecord) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)
	for _, record := range records {
		crop := strings.TrimSpace(record.MainCrop)
		if crop == "" {
			continue
		}
		if _, exists := seen[crop]; exists {
			continue
		}
		seen[crop] = struct{}{}
		result = append(result, crop)
	}
	sort.Strings(result)
	return result
}
