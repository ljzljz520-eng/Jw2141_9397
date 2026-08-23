package service

import (
	"fmt"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

func (r *Registry) FarmerTimeline(id string) ([]domain.AuditEvent, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("farmer id is required")
	}
	if _, err := r.GetFarmer(id); err != nil {
		return nil, err
	}
	return r.AuditTrail("FarmerRecord", id)
}

func (r *Registry) HasPublishedReport(farmerID string) (bool, error) {
	events, err := r.AuditTrail("VisitReport", farmerID)
	if err != nil {
		return false, err
	}
	for _, event := range events {
		if event.Action == domain.AuditPublished {
			return true, nil
		}
	}
	return false, nil
}

func (r *Registry) LatestWorkflowState(id string) (string, error) {
	event, err := r.LastAudit("FarmerRecord", id)
	if err != nil {
		return "", err
	}
	switch event.Action {
	case domain.AuditArchived:
		return domain.RecordArchived, nil
	case domain.AuditRestored:
		return domain.RecordActive, nil
	default:
		return domain.RecordActive, nil
	}
}
