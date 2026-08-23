package workflow

import (
	"fmt"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/service"
)

type CreateReviewArchive struct {
	Registry *service.Registry
}

func (w CreateReviewArchive) Run(record domain.FarmerRecord, report domain.VisitReport, reviewer string, archiveDate string) (domain.ArchiveEntry, error) {
	created, err := w.Registry.CreateFarmer(record)
	if err != nil {
		return domain.ArchiveEntry{}, err
	}
	report.FarmerID = created.ID
	savedReport, err := w.Registry.SubmitVisitReport(report)
	if err != nil {
		return domain.ArchiveEntry{}, err
	}
	review, err := w.Registry.OpenReview(savedReport, reviewer)
	if err != nil {
		return domain.ArchiveEntry{}, err
	}
	if _, err := w.Registry.ApproveReview(review.ID, "confirmed by township officer"); err != nil {
		return domain.ArchiveEntry{}, err
	}
	if _, err := w.Registry.PublishReport(savedReport.ID, review.ID); err != nil {
		return domain.ArchiveEntry{}, err
	}
	return w.Registry.ArchiveFarmer(created.ID, reviewer, archiveDate, "season closed")
}

func EnsureWorkflowReady(w CreateReviewArchive) error {
	if w.Registry == nil {
		return fmt.Errorf("workflow registry is required")
	}
	return nil
}
