package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"labops/model"
)

func CSVHeader() []string {
	return []string{"id", "equipment", "title", "status", "priority", "assignee", "updated_at"}
}

func CSVRow(record model.Record) []string {
	return []string{record.ID, record.EquipmentID, record.Title, record.Status, record.Priority, record.Assignee, record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")}
}

func ExportCSV(records []model.Record, output io.Writer) error {
	writer := csv.NewWriter(output)
	if err := writer.Write(CSVHeader()); err != nil {
		return err
	}
	for _, record := range records {
		if err := writer.Write(CSVRow(record)); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func ExportText(records []model.Record) string {
	var builder strings.Builder
	for _, record := range records {
		builder.WriteString(FormatRecord(record))
		builder.WriteByte('\n')
	}
	return builder.String()
}

func Detail(record model.Record, events []model.Event) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Record: %s\n", record.ID)
	fmt.Fprintf(&builder, "Equipment: %s\n", record.EquipmentID)
	fmt.Fprintf(&builder, "Title: %s\n", record.Title)
	fmt.Fprintf(&builder, "Status: %s\n", model.StatusLabel(record.Status))
	fmt.Fprintf(&builder, "Priority: %s\n", record.Priority)
	fmt.Fprintf(&builder, "Events: %d\n", len(events))
	return builder.String()
}

func ValidateExport(records []model.Record) error {
	for _, record := range records {
		if err := model.ValidateRecord(record); err != nil {
			return err
		}
	}
	return nil
}
