package validate

import (
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

func AuditEvent(event domain.AuditEvent) []string {
	errors := make([]string, 0, 5)
	if strings.TrimSpace(event.EntityType) == "" {
		errors = append(errors, "entity type is required")
	}
	if strings.TrimSpace(event.EntityID) == "" {
		errors = append(errors, "entity id is required")
	}
	if !domain.ValidAuditAction(event.Action) {
		errors = append(errors, "audit action is unsupported")
	}
	if strings.TrimSpace(event.Actor) == "" {
		errors = append(errors, "audit actor is required")
	}
	if !Date(event.OccurredOn) {
		errors = append(errors, "audit date must use YYYY-MM-DD")
	}
	return errors
}

func IsSystemAction(action string) bool {
	return action == domain.AuditImported || action == domain.AuditValidation
}

func FilterAudit(events []domain.AuditEvent, entityType string, entityID string) []domain.AuditEvent {
	result := make([]domain.AuditEvent, 0, len(events))
	for _, event := range events {
		if event.IsFor(entityType, entityID) {
			result = append(result, event)
		}
	}
	return result
}
