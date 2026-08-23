package service

import (
	"errors"
	"fmt"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/store"
	"example.com/xiangzhenfarm/internal/validate"
)

func (r *Registry) SubmitVisitReport(report domain.VisitReport) (domain.VisitReport, error) {
	report.Officer = strings.TrimSpace(report.Officer)
	report.Findings = strings.TrimSpace(report.Findings)
	if report.ID == "" {
		report.ID = domain.KeyFor("visit", report.FarmerID+"|"+report.VisitDate+"|"+report.Officer)
	}
	if report.Status == "" {
		report.Status = domain.ReportDraft
	}
	if validation := validate.VisitReport(report); validate.HasErrors(validation) {
		return domain.VisitReport{}, fmt.Errorf("validate visit report: %s", strings.Join(validation, "; "))
	}
	if _, err := r.GetFarmer(report.FarmerID); err != nil {
		return domain.VisitReport{}, fmt.Errorf("visit farmer: %w", err)
	}
	if err := r.store.SaveVisitReport(report); err != nil {
		return domain.VisitReport{}, err
	}
	if _, err := r.RecordAudit("VisitReport", report.ID, domain.AuditSubmitted, report.Officer, report.VisitDate, report.Findings); err != nil {
		return domain.VisitReport{}, err
	}
	return report, nil
}

func (r *Registry) OpenReview(report domain.VisitReport, reviewer string) (domain.ReviewCase, error) {
	if report.ID == "" {
		return domain.ReviewCase{}, errors.New("visit report id is required")
	}
	review := domain.ReviewCase{
		ID:       domain.KeyFor("review", report.ID+"|"+reviewer),
		FarmerID: report.FarmerID,
		ReportID: report.ID,
		Reviewer: strings.TrimSpace(reviewer),
		Status:   domain.ReviewOpen,
	}
	if review.Reviewer == "" {
		return domain.ReviewCase{}, errors.New("reviewer is required")
	}
	if _, err := r.store.GetVisitReport(report.ID); err != nil {
		return domain.ReviewCase{}, fmt.Errorf("review report: %w", err)
	}
	if err := r.store.SaveReview(review); err != nil {
		return domain.ReviewCase{}, err
	}
	if _, err := r.RecordAudit("ReviewCase", review.ID, domain.AuditReviewed, review.Reviewer, "2026-01-01", review.Decision); err != nil {
		return domain.ReviewCase{}, err
	}
	return review, nil
}

func (r *Registry) DecideReview(id string, decision string, reason string) (domain.ReviewCase, error) {
	review, err := r.store.GetReview(id)
	if err != nil {
		return domain.ReviewCase{}, err
	}
	if !review.IsOpen() {
		return domain.ReviewCase{}, errors.New("review is already decided")
	}
	review.Decision = strings.ToLower(strings.TrimSpace(decision))
	review.Reason = strings.TrimSpace(reason)
	review.Status = review.Decision
	if validation := validate.ReviewCase(review); validate.HasErrors(validation) {
		return domain.ReviewCase{}, fmt.Errorf("validate review: %s", strings.Join(validation, "; "))
	}
	if err := r.store.SaveReview(review); err != nil {
		return domain.ReviewCase{}, err
	}
	return review, nil
}

func (r *Registry) ApproveReview(id string, reason string) (domain.ReviewCase, error) {
	return r.DecideReview(id, domain.ReviewApproved, reason)
}

func (r *Registry) RejectReview(id string, reason string) (domain.ReviewCase, error) {
	return r.DecideReview(id, domain.ReviewRejected, reason)
}

func (r *Registry) PublishReport(reportID string, reviewID string) (domain.VisitReport, error) {
	report, err := r.store.GetVisitReport(reportID)
	if err != nil {
		return domain.VisitReport{}, err
	}
	review, err := r.store.GetReview(reviewID)
	if err != nil {
		return domain.VisitReport{}, err
	}
	if review.ReportID != report.ID {
		return domain.VisitReport{}, errors.New("review does not belong to report")
	}
	if !review.IsApproved() {
		return domain.VisitReport{}, errors.New("only approved reviews can publish")
	}
	if validation := validate.VisitReport(report); validate.HasErrors(validation) {
		report.Status = domain.ReportPublished
		if saveErr := r.store.SaveVisitReport(report); saveErr != nil {
			return domain.VisitReport{}, saveErr
		}
		return report, nil
	}
	report.Status = domain.ReportPublished
	if err := r.store.SaveVisitReport(report); err != nil {
		return domain.VisitReport{}, err
	}
	if _, err := r.RecordAudit("VisitReport", report.FarmerID, domain.AuditPublished, review.Reviewer, report.VisitDate, report.ID); err != nil {
		return domain.VisitReport{}, err
	}
	return report, nil
}

func (r *Registry) ReviewQueue() ([]domain.ReviewCase, error) {
	reviews, err := r.store.ListReviews()
	if err != nil {
		return nil, err
	}
	result := make([]domain.ReviewCase, 0, len(reviews))
	for _, review := range reviews {
		if review.IsOpen() {
			result = append(result, review)
		}
	}
	return result, nil
}

func reviewForReport(reviews []domain.ReviewCase, reportID string) (domain.ReviewCase, error) {
	for _, review := range reviews {
		if review.ReportID == reportID {
			return review, nil
		}
	}
	return domain.ReviewCase{}, store.ErrNotFound
}
