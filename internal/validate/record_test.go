package validate

import (
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
)

func TestFarmerRecordValidation(t *testing.T) {
	valid := domain.FarmerRecord{HouseholdHead: "Wang", VillageGroup: "North", CultivatedArea: 2.5, MainCrop: "rice", Phone: "13800000000", LastVisit: "2026-05-01"}
	if errors := FarmerRecord(valid); len(errors) != 0 {
		t.Fatalf("valid record rejected: %v", errors)
	}
	valid.Phone = ""
	if errors := FarmerRecord(valid); len(errors) == 0 {
		t.Fatal("invalid record accepted")
	}
}
