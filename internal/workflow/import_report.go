package workflow

import (
	"io"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/importer"
	"example.com/xiangzhenfarm/internal/service"
)

type ImportReport struct {
	Registry *service.Registry
}

func (w ImportReport) Run(files []importer.BatchFile) (importer.BatchResult, error) {
	return importer.ProcessFiles(files, w.Registry)
}

func SingleImport(name string, input io.Reader, registry *service.Registry) (importer.BatchResult, error) {
	return importer.ProcessFiles([]importer.BatchFile{{Name: name, Input: input}}, registry)
}

func RowReport(result importer.BatchResult, line int) (domain.ImportRow, bool) {
	for _, row := range result.Rows {
		if row.Line == line {
			return row, true
		}
	}
	return domain.ImportRow{}, false
}
