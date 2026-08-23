package store

import (
	"path/filepath"
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.db")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	farmer := domain.FarmerRecord{ID: "C-103", HouseholdHead: "Zhao", VillageGroup: "C", CultivatedArea: 1.2, MainCrop: "rice", Phone: "13800000000", LastVisit: "2026-06-01", Status: domain.RecordActive}
	report := domain.VisitReport{ID: "visit-C-103", FarmerID: farmer.ID, VisitDate: "2026-06-01", Officer: "Officer", Findings: "healthy", Status: domain.ReportPublished}
	review := domain.ReviewCase{ID: "review-C-103", FarmerID: farmer.ID, ReportID: report.ID, Reviewer: "Reviewer", Decision: domain.ReviewApproved, Status: domain.ReviewApproved}
	note := domain.CollaborationNote{ID: "note-C-103", FarmerID: farmer.ID, Author: "Worker", Body: "follow up", Visibility: "team"}
	archive := domain.ArchiveEntry{ID: "archive-C-103", FarmerID: farmer.ID, ArchivedBy: "Reviewer", ArchivedOn: "2026-06-02", Reason: "season close", Status: domain.ArchiveCurrent}
	for _, save := range []func() error{
		func() error { return database.SaveFarmer(farmer) },
		func() error { return database.SaveVisitReport(report) },
		func() error { return database.SaveReview(review) },
		func() error { return database.SaveNote(note) },
		func() error { return database.SaveArchive(archive) },
	} {
		if err := save(); err != nil {
			t.Fatal(err)
		}
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	loadedFarmer, err := reopened.GetFarmer(farmer.ID)
	if err != nil || loadedFarmer.HouseholdHead != farmer.HouseholdHead {
		t.Fatalf("farmer did not survive reopen: %#v %v", loadedFarmer, err)
	}
	loadedReport, err := reopened.GetVisitReport(report.ID)
	if err != nil || loadedReport.Findings != report.Findings {
		t.Fatalf("report did not survive reopen: %#v %v", loadedReport, err)
	}
	loadedReview, err := reopened.GetReview(review.ID)
	if err != nil || !loadedReview.IsApproved() {
		t.Fatalf("review did not survive reopen: %#v %v", loadedReview, err)
	}
	loadedNote, err := reopened.GetNote(note.ID)
	if err != nil || loadedNote.Body != note.Body {
		t.Fatalf("note did not survive reopen: %#v %v", loadedNote, err)
	}
	loadedArchive, err := reopened.GetArchive(archive.ID)
	if err != nil || !loadedArchive.IsCurrent() {
		t.Fatalf("archive did not survive reopen: %#v %v", loadedArchive, err)
	}
}
