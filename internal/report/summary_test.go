package report

import (
	"path/filepath"
	"strings"
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/service"
	"example.com/xiangzhenfarm/internal/store"
)

func TestSummaryAndExport(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "registry.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	registry := service.NewRegistry(database)
	record, err := registry.CreateFarmer(domain.FarmerRecord{HouseholdHead: "Qian", VillageGroup: "West", CultivatedArea: 4, MainCrop: "corn", Phone: "13800000003", LastVisit: "2026-05-03"})
	if err != nil {
		t.Fatal(err)
	}
	summary, err := Build(registry)
	if err != nil || summary.Total != 1 || summary.ByCrop["corn"] != 1 {
		t.Fatalf("summary failed: %#v %v", summary, err)
	}
	var output strings.Builder
	if err := WriteFarmers(&output, []domain.FarmerRecord{record}); err != nil || !strings.Contains(output.String(), "Qian") {
		t.Fatalf("export failed: %s %v", output.String(), err)
	}
}
