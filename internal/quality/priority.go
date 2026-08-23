package quality

import (
	"sort"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
)

type Priority struct {
	RecordID string
	Value    int
	Label    string
}

func VisitPriority(record domain.FarmerRecord) Priority {
	value := 1
	if record.CultivatedArea >= 10 {
		value += 3
	} else if record.CultivatedArea >= 5 {
		value += 2
	} else {
		value++
	}
	if strings.EqualFold(record.MainCrop, "rice") || strings.EqualFold(record.MainCrop, "corn") {
		value++
	}
	if record.Status == domain.RecordArchived {
		value = 0
	}
	return Priority{RecordID: record.ID, Value: value, Label: priorityLabel(value)}
}

func priorityLabel(value int) string {
	switch {
	case value >= 5:
		return "urgent"
	case value >= 3:
		return "normal"
	default:
		return "low"
	}
}

func Prioritize(records []domain.FarmerRecord) []Priority {
	result := make([]Priority, 0, len(records))
	for _, record := range records {
		result = append(result, VisitPriority(record))
	}
	sort.SliceStable(result, func(left, right int) bool {
		if result[left].Value == result[right].Value {
			return result[left].RecordID < result[right].RecordID
		}
		return result[left].Value > result[right].Value
	})
	return result
}

func PriorityMap(records []domain.FarmerRecord) map[string]Priority {
	result := make(map[string]Priority, len(records))
	for _, priority := range Prioritize(records) {
		result[priority.RecordID] = priority
	}
	return result
}

func IsUrgent(priority Priority) bool {
	return priority.Label == "urgent"
}

func CountUrgent(priorities []Priority) int {
	count := 0
	for _, priority := range priorities {
		if IsUrgent(priority) {
			count++
		}
	}
	return count
}

func CropPriority(record domain.FarmerRecord) string {
	if strings.TrimSpace(record.MainCrop) == "" {
		return "unknown"
	}
	return strings.ToLower(strings.TrimSpace(record.MainCrop))
}
