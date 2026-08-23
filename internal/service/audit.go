package service

import (
	"fmt"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/validate"
)

func (r *Registry) RecordAudit(entityType string, entityID string, action string, actor string, date string, details string) (domain.AuditEvent, error) {
	event := domain.AuditEvent{
		ID:         domain.KeyFor("audit", entityType+"|"+entityID+"|"+action+"|"+actor+"|"+date),
		EntityType: strings.TrimSpace(entityType),
		EntityID:   strings.TrimSpace(entityID),
		Action:     strings.TrimSpace(action),
		Actor:      strings.TrimSpace(actor),
		OccurredOn: strings.TrimSpace(date),
		Details:    strings.TrimSpace(details),
	}
	if validation := validate.AuditEvent(event); validate.HasErrors(validation) {
		return domain.AuditEvent{}, fmt.Errorf("validate audit event: %s", strings.Join(validation, "; "))
	}
	count, err := r.store.AuditCount()
	if err != nil {
		return domain.AuditEvent{}, err
	}
	event.Sequence = count + 1
	if err := r.store.SaveAuditEvent(event); err != nil {
		return domain.AuditEvent{}, err
	}
	return event, nil
}

func (r *Registry) AuditTrail(entityType string, entityID string) ([]domain.AuditEvent, error) {
	return r.store.AuditEventsFor(entityType, entityID)
}

func (r *Registry) AuditCount() (int, error) {
	return r.store.AuditCount()
}

func (r *Registry) LastAudit(entityType string, entityID string) (domain.AuditEvent, error) {
	events, err := r.AuditTrail(entityType, entityID)
	if err != nil {
		return domain.AuditEvent{}, err
	}
	if len(events) == 0 {
		return domain.AuditEvent{}, fmt.Errorf("no audit events for %s %s", entityType, entityID)
	}
	return events[len(events)-1], nil
}
