package importer

import (
	"strings"
	"testing"
)

func TestCSVImportSeparatesValidAndRejectedRows(t *testing.T) {
	input := "id,household_head,village_group,cultivated_area,main_crop,phone,last_visit\nC-1,Wang,North,2.5,rice,13800000000,2026-05-01\nC-2,Liu,South,0,wheat,13800000001,2026-05-02\n"
	result, err := NewCSVImporter().Read(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if result.Accepted != 1 || result.Rejected != 1 {
		t.Fatalf("unexpected import result: %#v", result)
	}
	if len(NewCSVImporter().ValidRows(result)) != 1 || len(NewCSVImporter().ErrorRows(result)) != 1 {
		t.Fatal("row partition is incorrect")
	}
}
