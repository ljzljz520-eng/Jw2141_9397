package domain

import "testing"

func TestRecordLifecycleStatus(t *testing.T) {
	record := FarmerRecord{ID: "farmer-1", HouseholdHead: "Li", VillageGroup: "East", Status: RecordActive}
	if !record.IsActive() || !record.CanEdit() || record.DisplayName() != "Li / East" {
		t.Fatal("active record state is incorrect")
	}
	record.Status = RecordArchived
	if !record.IsArchived() || record.CanEdit() {
		t.Fatal("archived record state is incorrect")
	}
}
