package service

import (
	"fmt"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/validate"
)

func (r *Registry) ArchiveFarmer(id string, operator string, date string, reason string) (domain.ArchiveEntry, error) {
	record, err := r.GetFarmer(id)
	if err != nil {
		return domain.ArchiveEntry{}, err
	}
	if record.IsArchived() {
		return domain.ArchiveEntry{}, fmt.Errorf("farmer already archived: %s", id)
	}
	entry := domain.ArchiveEntry{
		ID:         domain.KeyFor("archive", id+"|"+date),
		FarmerID:   id,
		ArchivedBy: strings.TrimSpace(operator),
		ArchivedOn: strings.TrimSpace(date),
		Reason:     strings.TrimSpace(reason),
		Status:     domain.ArchiveCurrent,
	}
	if validation := validate.ArchiveEntry(entry); validate.HasErrors(validation) {
		return domain.ArchiveEntry{}, fmt.Errorf("validate archive: %s", strings.Join(validation, "; "))
	}
	record.Status = domain.RecordArchived
	if err := r.store.SaveFarmer(record); err != nil {
		return domain.ArchiveEntry{}, err
	}
	if err := r.store.SaveArchive(entry); err != nil {
		return domain.ArchiveEntry{}, err
	}
	if _, err := r.RecordAudit("FarmerRecord", id, domain.AuditArchived, operator, date, reason); err != nil {
		return domain.ArchiveEntry{}, err
	}
	return entry, nil
}

func (r *Registry) RestoreFarmer(id string) (domain.FarmerRecord, error) {
	record, err := r.GetFarmer(id)
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	if !record.IsArchived() {
		return domain.FarmerRecord{}, fmt.Errorf("farmer is not archived: %s", id)
	}
	record.Status = domain.RecordActive
	if err := r.store.SaveFarmer(record); err != nil {
		return domain.FarmerRecord{}, err
	}
	if _, err := r.RecordAudit("FarmerRecord", id, domain.AuditRestored, "registry", record.LastVisit, "restored from archive"); err != nil {
		return domain.FarmerRecord{}, err
	}
	return record, nil
}

func (r *Registry) ArchiveHistory(id string) ([]domain.ArchiveEntry, error) {
	return r.store.ArchivesForFarmer(id)
}

func (r *Registry) ArchiveCount() (int, error) {
	entries, err := r.store.ListArchives()
	return len(entries), err
}
