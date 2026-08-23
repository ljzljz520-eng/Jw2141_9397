package workflow

import (
	"path/filepath"
	"strings"
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/importer"
	"example.com/xiangzhenfarm/internal/service"
	"example.com/xiangzhenfarm/internal/store"
)

func newWorkflowRegistry(t *testing.T) *service.Registry {
	t.Helper()
	database, err := store.Open(filepath.Join(t.TempDir(), "registry.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	return service.NewRegistry(database)
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	registry := newWorkflowRegistry(t)
	workflow := CreateReviewArchive{Registry: registry}
	entry, err := workflow.Run(domain.FarmerRecord{HouseholdHead: "Sun", VillageGroup: "River", CultivatedArea: 2, MainCrop: "rice", Phone: "13800000004", LastVisit: "2026-05-04"}, domain.VisitReport{VisitDate: "2026-05-04", Officer: "Officer", Findings: "ready"}, "Officer", "2026-05-05")
	if err != nil || entry.Status != domain.ArchiveCurrent {
		t.Fatalf("workflow failed: %#v %v", entry, err)
	}
	record, err := registry.GetFarmer(entry.FarmerID)
	if err != nil || !record.IsArchived() {
		t.Fatalf("archive state missing: %#v %v", record, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	registry := newWorkflowRegistry(t)
	record, err := registry.CreateFarmer(domain.FarmerRecord{HouseholdHead: "He", VillageGroup: "Hill", CultivatedArea: 2, MainCrop: "tea", Phone: "13800000005", LastVisit: "2026-05-05"})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := (SearchUpdatePublish{Registry: registry}).Run("tea", domain.CollaborationNote{FarmerID: record.ID, Author: "Worker", Body: "new visit", Visibility: "team"}, record.ID)
	if err != nil || updated.Notes != "new visit" {
		t.Fatalf("workflow failed: %#v %v", updated, err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	registry := newWorkflowRegistry(t)
	input := "id,household_head,village_group,cultivated_area,main_crop,phone,last_visit\nC-201,Lin,Plain,3,rice,13800000006,2026-05-06\n"
	result, err := (ImportReport{Registry: registry}).Run([]importer.BatchFile{{Name: "C-201", Input: strings.NewReader(input)}})
	if err != nil || result.Imported != 1 || result.Rejected != 0 {
		t.Fatalf("workflow failed: %#v %v", result, err)
	}
	record, err := registry.GetFarmer("C-201")
	if err != nil || record.MainCrop != "rice" {
		t.Fatalf("imported record missing: %#v %v", record, err)
	}
}
