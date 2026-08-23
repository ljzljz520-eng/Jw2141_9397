package domain

func (r FarmerRecord) IsActive() bool {
	return r.Status == RecordActive || r.Status == ""
}

func (r FarmerRecord) IsArchived() bool {
	return r.Status == RecordArchived
}

func (r FarmerRecord) CanEdit() bool {
	if r.IsArchived() {
		return false
	}
	return r.ID != ""
}

func (r FarmerRecord) DisplayName() string {
	if r.HouseholdHead == "" {
		return r.ID
	}
	return r.HouseholdHead + " / " + r.VillageGroup
}

func (r VisitReport) IsDraft() bool {
	return r.Status == ReportDraft || r.Status == ""
}

func (r VisitReport) IsPublished() bool {
	return r.Status == ReportPublished
}

func (r ReviewCase) IsOpen() bool {
	return r.Status == ReviewOpen || r.Status == ""
}

func (r ReviewCase) IsApproved() bool {
	return r.Status == ReviewApproved
}

func (a ArchiveEntry) IsCurrent() bool {
	return a.Status == ArchiveCurrent || a.Status == ""
}

func NormalizeStatus(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
