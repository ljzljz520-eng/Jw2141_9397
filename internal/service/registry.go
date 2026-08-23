package service

import (
	"errors"
	"fmt"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/store"
	"example.com/xiangzhenfarm/internal/validate"
)

type Registry struct {
	store *store.Database
}

func NewRegistry(database *store.Database) *Registry {
	return &Registry{store: database}
}

func (r *Registry) CreateFarmer(record domain.FarmerRecord) (domain.FarmerRecord, error) {
	record.HouseholdHead = strings.TrimSpace(record.HouseholdHead)
	record.VillageGroup = strings.TrimSpace(record.VillageGroup)
	record.MainCrop = strings.TrimSpace(record.MainCrop)
	record.Phone = strings.TrimSpace(record.Phone)
	record.LastVisit = strings.TrimSpace(record.LastVisit)
	if record.ID == "" {
		record.ID = domain.KeyFor("farmer", record.HouseholdHead+"|"+record.VillageGroup+"|"+record.MainCrop)
	}
	if record.Status == "" {
		record.Status = domain.RecordActive
	}
	if validation := validate.FarmerRecord(record); validate.HasErrors(validation) {
		return domain.FarmerRecord{}, fmt.Errorf("validate farmer: %s", strings.Join(validation, "; "))
	}
	if _, err := r.store.GetFarmer(record.ID); err == nil {
		return domain.FarmerRecord{}, fmt.Errorf("farmer already exists: %s", record.ID)
	}
	if err := r.store.SaveFarmer(record); err != nil {
		return domain.FarmerRecord{}, fmt.Errorf("save farmer: %w", err)
	}
	if _, err := r.RecordAudit("FarmerRecord", record.ID, domain.AuditCreated, "registry", record.LastVisit, "manual registration"); err != nil {
		return domain.FarmerRecord{}, err
	}
	return record, nil
}

func (r *Registry) GetFarmer(id string) (domain.FarmerRecord, error) {
	if strings.TrimSpace(id) == "" {
		return domain.FarmerRecord{}, errors.New("farmer id is required")
	}
	return r.store.GetFarmer(id)
}

func (r *Registry) UpdateFarmer(id string, changes domain.FarmerRecord) (domain.FarmerRecord, error) {
	record, err := r.GetFarmer(id)
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	if !record.CanEdit() {
		return domain.FarmerRecord{}, errors.New("archived farmer cannot be edited")
	}
	updated := mergeFarmer(record, changes)
	if validation := validate.FarmerRecord(updated); validate.HasErrors(validation) {
		return domain.FarmerRecord{}, fmt.Errorf("validate farmer update: %s", strings.Join(validation, "; "))
	}
	if err := r.store.SaveFarmer(updated); err != nil {
		return domain.FarmerRecord{}, fmt.Errorf("save farmer update: %w", err)
	}
	if _, err := r.RecordAudit("FarmerRecord", updated.ID, domain.AuditUpdated, "registry", updated.LastVisit, "record change"); err != nil {
		return domain.FarmerRecord{}, err
	}
	return updated, nil
}

func mergeFarmer(original domain.FarmerRecord, changes domain.FarmerRecord) domain.FarmerRecord {
	if changes.HouseholdHead != "" {
		original.HouseholdHead = strings.TrimSpace(changes.HouseholdHead)
	}
	if changes.VillageGroup != "" {
		original.VillageGroup = strings.TrimSpace(changes.VillageGroup)
	}
	if changes.CultivatedArea > 0 {
		original.CultivatedArea = changes.CultivatedArea
	}
	if changes.MainCrop != "" {
		original.MainCrop = strings.TrimSpace(changes.MainCrop)
	}
	if changes.Phone != "" {
		original.Phone = strings.TrimSpace(changes.Phone)
	}
	if changes.LastVisit != "" {
		original.LastVisit = strings.TrimSpace(changes.LastVisit)
	}
	if changes.Notes != "" {
		original.Notes = strings.TrimSpace(changes.Notes)
	}
	return original
}

func (r *Registry) DeleteFarmer(id string) error {
	record, err := r.GetFarmer(id)
	if err != nil {
		return err
	}
	if record.IsArchived() {
		return errors.New("archived farmer must be restored before deletion")
	}
	if err := r.store.DeleteFarmer(id); err != nil {
		return err
	}
	return nil
}

func (r *Registry) AllFarmers() ([]domain.FarmerRecord, error) {
	return r.store.ListFarmers()
}

func (r *Registry) SearchFarmers(filter domain.SearchFilter, sortSpec domain.SortSpec) ([]domain.FarmerRecord, error) {
	records, err := r.AllFarmers()
	if err != nil {
		return nil, err
	}
	filtered := validate.FilterRecords(records, filter)
	return validate.SortRecords(filtered, sortSpec), nil
}

func (r *Registry) SearchByCrop(crop string) ([]domain.FarmerRecord, error) {
	return r.SearchFarmers(domain.SearchFilter{MainCrop: crop}, domain.SortSpec{Field: "village_group"})
}

func (r *Registry) Count() (int, error) {
	return r.store.FarmerCount()
}

func (r *Registry) Store() *store.Database {
	return r.store
}
