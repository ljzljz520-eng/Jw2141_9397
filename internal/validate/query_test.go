package validate

import (
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
)

func TestCropFilterAndAreaSort(t *testing.T) {
	records := []domain.FarmerRecord{
		{ID: "1", HouseholdHead: "A", VillageGroup: "North", CultivatedArea: 3, MainCrop: "rice"},
		{ID: "2", HouseholdHead: "B", VillageGroup: "South", CultivatedArea: 1, MainCrop: "wheat"},
		{ID: "3", HouseholdHead: "C", VillageGroup: "North", CultivatedArea: 2, MainCrop: "rice"},
	}
	filtered := FilterRecords(records, domain.SearchFilter{MainCrop: "RICE", VillageGroup: "north"})
	if len(filtered) != 2 {
		t.Fatalf("expected two rice records, got %d", len(filtered))
	}
	sorted := SortRecords(filtered, domain.SortSpec{Field: "area"})
	if sorted[0].ID != "3" {
		t.Fatalf("unexpected sort order: %#v", sorted)
	}
}
