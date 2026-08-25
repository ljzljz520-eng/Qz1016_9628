package service

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"labops/model"
)

type ImportResult struct {
	Accepted int
	Rejected int
	Errors   []string
}

func ParseRecordRow(row []string, now time.Time) (model.Record, error) {
	if len(row) < 3 {
		return model.Record{}, fmt.Errorf("record row needs three columns")
	}
	record := model.NewRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[1]), strings.TrimSpace(row[2]), now)
	if len(row) > 3 {
		record.Description = strings.TrimSpace(row[3])
	}
	if len(row) > 4 {
		record.Priority = model.NormalizePriority(strings.TrimSpace(row[4]))
	}
	if len(row) > 5 {
		record.Assignee = strings.TrimSpace(row[5])
	}
	if err := model.ValidateRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func ImportCSV(reader io.Reader, now time.Time, save func(model.Record) error) ImportResult {
	result := ImportResult{Errors: make([]string, 0)}
	csvReader := csv.NewReader(reader)
	line := 0
	for {
		line++
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		if line == 1 && strings.EqualFold(strings.TrimSpace(row[0]), "id") {
			continue
		}
		record, parseErr := ParseRecordRow(row, now)
		if parseErr != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, parseErr))
			continue
		}
		if saveErr := save(record); saveErr != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, saveErr))
			continue
		}
		result.Accepted++
	}
	return result
}

func ImportLines(reader io.Reader, now time.Time, save func(model.Record) error) ImportResult {
	result := ImportResult{Errors: make([]string, 0)}
	scanner := bufio.NewScanner(reader)
	line := 0
	for scanner.Scan() {
		line++
		fields := strings.Split(scanner.Text(), "|")
		record, err := ParseRecordRow(fields, now)
		if err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		if err = save(record); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("line %d: %v", line, err))
			continue
		}
		result.Accepted++
	}
	if err := scanner.Err(); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	return result
}

func NormalizeImport(result ImportResult) ImportResult {
	if result.Accepted < 0 {
		result.Accepted = 0
	}
	if result.Rejected < 0 {
		result.Rejected = 0
	}
	if result.Errors == nil {
		result.Errors = []string{}
	}
	return result
}
