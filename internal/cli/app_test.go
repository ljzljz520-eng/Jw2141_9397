package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"example.com/xiangzhenfarm/internal/service"
	"example.com/xiangzhenfarm/internal/store"
)

func TestCommandLineAddAndSearch(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "registry.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	var output bytes.Buffer
	app := New(service.NewRegistry(database), strings.NewReader(""), &output)
	if err := app.Run([]string{"add", "-head", "Mo", "-village", "Lake", "-area", "1.5", "-crop", "rice", "-phone", "13800000007", "-visit", "2026-05-07"}); err != nil {
		t.Fatal(err)
	}
	output.Reset()
	if err := app.Run([]string{"search", "-crop", "rice"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Mo") {
		t.Fatalf("search output missing record: %s", output.String())
	}
}
