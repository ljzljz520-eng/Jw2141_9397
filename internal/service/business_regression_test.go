package service

import (
	"path/filepath"
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/store"
)

func TestBusiness003Regression(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "registry.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	registry := NewRegistry(database)
	farmer, err := registry.CreateFarmer(domain.FarmerRecord{ID: "C-103", HouseholdHead: "Zhao", VillageGroup: "C", CultivatedArea: 1.2, MainCrop: "rice", Phone: "13800000000", LastVisit: "2026-06-01"})
	if err != nil {
		t.Fatal(err)
	}
	report := domain.VisitReport{ID: "report-C-103", FarmerID: farmer.ID, VisitDate: "2026-06-01", Officer: "Officer", Findings: "", MissingField: "findings", Status: domain.ReportDraft, Sequence: 25}
	if err := registry.Store().SaveVisitReport(report); err != nil {
		t.Fatal(err)
	}
	review, err := registry.OpenReview(report, "Reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.ApproveReview(review.ID, "reviewed"); err != nil {
		t.Fatal(err)
	}
	published, err := registry.PublishReport(report.ID, review.ID)
	if err == nil {
		t.Fatalf("invalid report was published: %#v", published)
	}
	if published.Findings != "" {
		t.Fatalf("missing value changed: %#v", published)
	}
	stored, err := registry.Store().GetVisitReport(report.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status == domain.ReportPublished {
		t.Fatalf("invalid report persisted as published: %#v", stored)
	}
}
