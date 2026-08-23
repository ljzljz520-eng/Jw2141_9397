package service

import (
	"path/filepath"
	"testing"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/store"
)

func TestRegistryCrudAndCropSearch(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "registry.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	registry := NewRegistry(database)
	one, err := registry.CreateFarmer(domain.FarmerRecord{HouseholdHead: "Wang", VillageGroup: "North", CultivatedArea: 3, MainCrop: "rice", Phone: "13800000001", LastVisit: "2026-05-01"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.CreateFarmer(domain.FarmerRecord{HouseholdHead: "Liu", VillageGroup: "South", CultivatedArea: 1, MainCrop: "wheat", Phone: "13800000002", LastVisit: "2026-05-02"}); err != nil {
		t.Fatal(err)
	}
	found, err := registry.SearchByCrop("rice")
	if err != nil || len(found) != 1 || found[0].ID != one.ID {
		t.Fatalf("crop search failed: %#v %v", found, err)
	}
	updated, err := registry.UpdateFarmer(one.ID, domain.FarmerRecord{Phone: "13800000009"})
	if err != nil || updated.Phone != "13800000009" {
		t.Fatalf("update failed: %#v %v", updated, err)
	}
	if err := registry.DeleteFarmer(one.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.GetFarmer(one.ID); err == nil {
		t.Fatal("deleted farmer still exists")
	}
}
