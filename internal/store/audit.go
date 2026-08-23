package store

import (
	"sort"

	"example.com/xiangzhenfarm/internal/domain"
)

const auditBucket = "audit_events"

func init() {
	bucketNames = append(bucketNames, []byte(auditBucket))
}

func (s *Database) SaveAuditEvent(event domain.AuditEvent) error {
	return s.Put(auditBucket, event.ID, event)
}

func (s *Database) GetAuditEvent(id string) (domain.AuditEvent, error) {
	var event domain.AuditEvent
	err := s.Get(auditBucket, id, &event)
	return event, err
}

func (s *Database) ListAuditEvents() ([]domain.AuditEvent, error) {
	keys, err := s.Keys(auditBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEvent, 0, len(keys))
	for _, key := range keys {
		event, getErr := s.GetAuditEvent(key)
		if getErr != nil {
			return nil, getErr
		}
		result = append(result, event)
	}
	sort.SliceStable(result, func(left, right int) bool {
		if result[left].Sequence == result[right].Sequence {
			return result[left].ID < result[right].ID
		}
		return result[left].Sequence < result[right].Sequence
	})
	return result, nil
}

func (s *Database) AuditEventsFor(entityType string, entityID string) ([]domain.AuditEvent, error) {
	all, err := s.ListAuditEvents()
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEvent, 0)
	for _, event := range all {
		if event.IsFor(entityType, entityID) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Database) AuditCount() (int, error) {
	return s.Count(auditBucket)
}
