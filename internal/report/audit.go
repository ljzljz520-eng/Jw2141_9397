package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

func RenderTimeline(events []domain.AuditEvent) string {
	lines := make([]string, 0, len(events))
	for _, event := range events {
		lines = append(lines, fmt.Sprintf("%03d %s %s %s %s", event.Sequence, event.OccurredOn, event.EntityType, event.Label(), event.Details))
	}
	return strings.Join(lines, "\n")
}

func WriteAudit(output io.Writer, events []domain.AuditEvent) error {
	writer := csv.NewWriter(output)
	if err := writer.Write([]string{"sequence", "entity_type", "entity_id", "action", "actor", "occurred_on", "details"}); err != nil {
		return err
	}
	for _, event := range events {
		if err := writer.Write([]string{fmt.Sprint(event.Sequence), event.EntityType, event.EntityID, event.Action, event.Actor, event.OccurredOn, event.Details}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func CountActions(events []domain.AuditEvent) map[string]int {
	result := make(map[string]int)
	for _, event := range events {
		result[event.Action]++
	}
	return result
}

func HasAction(events []domain.AuditEvent, action string) bool {
	for _, event := range events {
		if event.Action == action {
			return true
		}
	}
	return false
}
