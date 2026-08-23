package domain

const (
	RecordActive    = "active"
	RecordArchived  = "archived"
	ReviewOpen      = "open"
	ReviewApproved  = "approved"
	ReviewRejected  = "rejected"
	ReportDraft     = "draft"
	ReportPublished = "published"
	ArchiveCurrent  = "current"
	ArchiveRestored = "restored"
)

type FarmerRecord struct {
	ID             string  `json:"id"`
	HouseholdHead  string  `json:"household_head"`
	VillageGroup   string  `json:"village_group"`
	CultivatedArea float64 `json:"cultivated_area"`
	MainCrop       string  `json:"main_crop"`
	Phone          string  `json:"phone"`
	LastVisit      string  `json:"last_visit"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes"`
	CreatedFrom    string  `json:"created_from"`
}

type VisitReport struct {
	ID           string `json:"id"`
	FarmerID     string `json:"farmer_id"`
	VisitDate    string `json:"visit_date"`
	Officer      string `json:"officer"`
	Findings     string `json:"findings"`
	MissingField string `json:"missing_field"`
	Status       string `json:"status"`
	Sequence     int    `json:"sequence"`
}

type ReviewCase struct {
	ID       string `json:"id"`
	FarmerID string `json:"farmer_id"`
	ReportID string `json:"report_id"`
	Reviewer string `json:"reviewer"`
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
	Status   string `json:"status"`
}

type CollaborationNote struct {
	ID         string `json:"id"`
	FarmerID   string `json:"farmer_id"`
	Author     string `json:"author"`
	Body       string `json:"body"`
	CreatedAt  string `json:"created_at"`
	Visibility string `json:"visibility"`
}

type ArchiveEntry struct {
	ID         string `json:"id"`
	FarmerID   string `json:"farmer_id"`
	ArchivedBy string `json:"archived_by"`
	ArchivedOn string `json:"archived_on"`
	Reason     string `json:"reason"`
	Status     string `json:"status"`
}

type SearchFilter struct {
	HouseholdHead string
	VillageGroup  string
	MainCrop      string
	Status        string
	MinArea       float64
	MaxArea       float64
}

type SortSpec struct {
	Field   string
	Reverse bool
}

type ImportRow struct {
	Line   int
	Record FarmerRecord
	Errors []string
}

type ImportSummary struct {
	Accepted int
	Rejected int
	Rows     []ImportRow
}
