package report

import (
	"fmt"
	"sort"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/quality"
	"example.com/xiangzhenfarm/internal/service"
)

type Summary struct {
	Total      int
	Active     int
	Archived   int
	ByVillage  map[string]int
	ByCrop     map[string]int
	Quality    []quality.Assessment
	Priority   map[string]quality.Priority
	Priorities []quality.Priority
}

func Build(registry *service.Registry) (Summary, error) {
	records, err := registry.AllFarmers()
	if err != nil {
		return Summary{}, err
	}
	summary := Summary{ByVillage: make(map[string]int), ByCrop: make(map[string]int), Quality: quality.AssessAll(records), Priority: quality.PriorityMap(records), Priorities: quality.Prioritize(records)}
	for _, record := range records {
		summary.Total++
		if record.IsArchived() {
			summary.Archived++
		} else {
			summary.Active++
		}
		summary.ByVillage[record.VillageGroup]++
		summary.ByCrop[record.MainCrop]++
	}
	return summary, nil
}

func Render(summary Summary) string {
	lines := []string{fmt.Sprintf("total=%d active=%d archived=%d", summary.Total, summary.Active, summary.Archived)}
	lines = append(lines, renderMap("village", summary.ByVillage)...)
	lines = append(lines, renderMap("crop", summary.ByCrop)...)
	lines = append(lines, fmt.Sprintf("quality_assessments=%d urgent_visits=%d", len(summary.Quality), quality.CountUrgent(summary.Priorities)))
	return strings.Join(lines, "\n")
}

func renderMap(prefix string, values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("%s=%s:%d", prefix, key, values[key]))
	}
	return lines
}

func RecordLine(record domain.FarmerRecord) string {
	return fmt.Sprintf("%s\t%s\t%s\t%.2f\t%s\t%s\t%s\t%s", record.ID, record.HouseholdHead, record.VillageGroup, record.CultivatedArea, record.MainCrop, record.Phone, record.LastVisit, record.Status)
}
