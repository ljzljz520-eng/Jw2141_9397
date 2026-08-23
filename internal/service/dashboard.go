package service

import (
	"fmt"
	"sort"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/quality"
)

type Dashboard struct {
	TotalFarmers      int
	ActiveFarmers     int
	ArchivedFarmers   int
	OpenReviews       int
	StoredNotes       int
	StoredArchives    int
	AuditEvents       int
	UrgentVisits      int
	IncompleteRecords int
	CropCounts        map[string]int
	VillageCounts     map[string]int
	Priorities        []quality.Priority
}

func (r *Registry) BuildDashboard() (Dashboard, error) {
	records, err := r.AllFarmers()
	if err != nil {
		return Dashboard{}, err
	}
	reviews, err := r.ReviewQueue()
	if err != nil {
		return Dashboard{}, err
	}
	notes, err := r.store.ListNotes()
	if err != nil {
		return Dashboard{}, err
	}
	archives, err := r.store.ListArchives()
	if err != nil {
		return Dashboard{}, err
	}
	auditEvents, err := r.store.ListAuditEvents()
	if err != nil {
		return Dashboard{}, err
	}
	dashboard := Dashboard{
		TotalFarmers:   len(records),
		OpenReviews:    len(reviews),
		StoredNotes:    len(notes),
		StoredArchives: len(archives),
		AuditEvents:    len(auditEvents),
		CropCounts:     make(map[string]int),
		VillageCounts:  make(map[string]int),
		Priorities:     quality.Prioritize(records),
	}
	for _, record := range records {
		if record.IsArchived() {
			dashboard.ArchivedFarmers++
		} else {
			dashboard.ActiveFarmers++
		}
		dashboard.CropCounts[record.MainCrop]++
		dashboard.VillageCounts[record.VillageGroup]++
		assessment := quality.Assess(record)
		if !quality.CanPublish(assessment) {
			dashboard.IncompleteRecords++
		}
		if quality.IsUrgent(quality.VisitPriority(record)) {
			dashboard.UrgentVisits++
		}
	}
	return dashboard, nil
}

func (dashboard Dashboard) IsHealthy() bool {
	if dashboard.TotalFarmers == 0 {
		return false
	}
	return dashboard.IncompleteRecords == 0 && dashboard.OpenReviews == 0
}

func (dashboard Dashboard) CompletionRate() float64 {
	if dashboard.TotalFarmers == 0 {
		return 0
	}
	return float64(dashboard.TotalFarmers-dashboard.IncompleteRecords) * 100 / float64(dashboard.TotalFarmers)
}

func (dashboard Dashboard) CropNames() []string {
	return sortedKeys(dashboard.CropCounts)
}

func (dashboard Dashboard) VillageNames() []string {
	return sortedKeys(dashboard.VillageCounts)
}

func (dashboard Dashboard) CropCount(crop string) int {
	return dashboard.CropCounts[strings.TrimSpace(crop)]
}

func (dashboard Dashboard) VillageCount(village string) int {
	return dashboard.VillageCounts[strings.TrimSpace(village)]
}

func (dashboard Dashboard) PriorityFor(id string) (quality.Priority, bool) {
	for _, priority := range dashboard.Priorities {
		if priority.RecordID == id {
			return priority, true
		}
	}
	return quality.Priority{}, false
}

func sortedKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (dashboard Dashboard) Render() string {
	lines := []string{
		fmt.Sprintf("farmers total=%d active=%d archived=%d", dashboard.TotalFarmers, dashboard.ActiveFarmers, dashboard.ArchivedFarmers),
		fmt.Sprintf("workflow reviews_open=%d notes=%d archives=%d audit_events=%d", dashboard.OpenReviews, dashboard.StoredNotes, dashboard.StoredArchives, dashboard.AuditEvents),
		fmt.Sprintf("quality incomplete=%d completion=%.1f%% urgent_visits=%d healthy=%t", dashboard.IncompleteRecords, dashboard.CompletionRate(), dashboard.UrgentVisits, dashboard.IsHealthy()),
	}
	for _, crop := range dashboard.CropNames() {
		lines = append(lines, fmt.Sprintf("crop %s=%d", crop, dashboard.CropCounts[crop]))
	}
	for _, village := range dashboard.VillageNames() {
		lines = append(lines, fmt.Sprintf("village %s=%d", village, dashboard.VillageCounts[village]))
	}
	return strings.Join(lines, "\n")
}

func (r *Registry) NeedsFollowUp() ([]domain.FarmerRecord, error) {
	records, err := r.AllFarmers()
	if err != nil {
		return nil, err
	}
	result := make([]domain.FarmerRecord, 0)
	for _, record := range records {
		if record.IsArchived() {
			continue
		}
		if !quality.CanPublish(quality.Assess(record)) {
			result = append(result, record)
			continue
		}
		if record.Notes == "" {
			result = append(result, record)
		}
	}
	return result, nil
}

func (r *Registry) RecordsByVillage(village string) ([]domain.FarmerRecord, error) {
	return r.SearchFarmers(domain.SearchFilter{VillageGroup: village}, domain.SortSpec{Field: "household_head"})
}

func (r *Registry) RecordsByCrop(crop string) ([]domain.FarmerRecord, error) {
	return r.SearchFarmers(domain.SearchFilter{MainCrop: crop}, domain.SortSpec{Field: "village_group"})
}

func (r *Registry) ReviewHealth() (string, error) {
	queue, err := r.ReviewQueue()
	if err != nil {
		return "", err
	}
	if len(queue) == 0 {
		return "clear", nil
	}
	if len(queue) > 10 {
		return "overloaded", nil
	}
	return "pending", nil
}

func (r *Registry) ArchiveHealth() (string, error) {
	count, err := r.ArchiveCount()
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "empty", nil
	}
	return "tracked", nil
}
