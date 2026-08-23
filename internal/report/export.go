package report

import (
	"encoding/csv"
	"io"

	"example.com/xiangzhenfarm/internal/domain"
)

func WriteFarmers(output io.Writer, records []domain.FarmerRecord) error {
	writer := csv.NewWriter(output)
	if err := writer.Write([]string{"id", "household_head", "village_group", "cultivated_area", "main_crop", "phone", "last_visit", "status"}); err != nil {
		return err
	}
	for _, record := range records {
		if err := writer.Write([]string{record.ID, record.HouseholdHead, record.VillageGroup, formatArea(record.CultivatedArea), record.MainCrop, record.Phone, record.LastVisit, record.Status}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func formatArea(area float64) string {
	return formatFloat(area)
}
