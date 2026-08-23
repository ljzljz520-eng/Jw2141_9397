package workflow

import (
	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/service"
)

type SearchUpdatePublish struct {
	Registry *service.Registry
}

func (w SearchUpdatePublish) Run(crop string, note domain.CollaborationNote, id string) (domain.FarmerRecord, error) {
	records, err := w.Registry.SearchByCrop(crop)
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	if len(records) == 0 {
		return domain.FarmerRecord{}, service.ErrNoSearchResults
	}
	selected, err := w.Registry.GetFarmer(id)
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	if selected.MainCrop != records[0].MainCrop && selected.ID != records[0].ID {
		return domain.FarmerRecord{}, service.ErrNoSearchResults
	}
	selected, err = w.Registry.UpdateFarmer(id, domain.FarmerRecord{Notes: note.Body})
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	if _, err := w.Registry.AddCollaborationNote(note); err != nil {
		return domain.FarmerRecord{}, err
	}
	report, err := w.Registry.SubmitVisitReport(domain.VisitReport{FarmerID: selected.ID, VisitDate: selected.LastVisit, Officer: note.Author, Findings: "field update confirmed"})
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	review, err := w.Registry.OpenReview(report, note.Author)
	if err != nil {
		return domain.FarmerRecord{}, err
	}
	if _, err := w.Registry.ApproveReview(review.ID, "published after field update"); err != nil {
		return domain.FarmerRecord{}, err
	}
	if _, err := w.Registry.PublishReport(report.ID, review.ID); err != nil {
		return domain.FarmerRecord{}, err
	}
	return selected, nil
}

func FirstMatch(records []domain.FarmerRecord) (domain.FarmerRecord, bool) {
	if len(records) == 0 {
		return domain.FarmerRecord{}, false
	}
	return records[0], true
}
