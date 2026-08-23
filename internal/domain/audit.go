package domain

const (
	AuditCreated    = "created"
	AuditUpdated    = "updated"
	AuditImported   = "imported"
	AuditSubmitted  = "submitted"
	AuditReviewed   = "reviewed"
	AuditPublished  = "published"
	AuditArchived   = "archived"
	AuditRestored   = "restored"
	AuditNoteAdded  = "note_added"
	AuditValidation = "validation_failed"
)

type AuditEvent struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	OccurredOn string `json:"occurred_on"`
	Details    string `json:"details"`
	Sequence   int    `json:"sequence"`
}

func (event AuditEvent) IsFor(entityType string, entityID string) bool {
	return event.EntityType == entityType && event.EntityID == entityID
}

func (event AuditEvent) Label() string {
	if event.Actor == "" {
		return event.Action
	}
	return event.Action + " by " + event.Actor
}

func ValidAuditAction(action string) bool {
	switch action {
	case AuditCreated, AuditUpdated, AuditImported, AuditSubmitted, AuditReviewed, AuditPublished, AuditArchived, AuditRestored, AuditNoteAdded, AuditValidation:
		return true
	default:
		return false
	}
}
