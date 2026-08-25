package report

import (
	"fmt"
	"labops/model"
	"strings"
)

func Bar(value, total, width int) string {
	if width < 1 {
		return ""
	}
	if total <= 0 {
		total = 1
	}
	filled := value * width / total
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("#", filled) + strings.Repeat("-", width-filled)
}
func StatusChart(counts map[string]int, width int) string {
	total := 0
	for _, n := range counts {
		total += n
	}
	keys := []string{"new", "validated", "processing", "blocked", "resolved", "closed", "archived", "cancelled"}
	lines := make([]string, 0)
	for _, key := range keys {
		if n := counts[key]; n > 0 {
			lines = append(lines, fmt.Sprintf("%-10s %s %d", key, Bar(n, total, width), n))
		}
	}
	return strings.Join(lines, "\n")
}
func PriorityChart(counts map[string]int, width int) string {
	total := 0
	for _, n := range counts {
		total += n
	}
	keys := []string{"urgent", "high", "normal", "low"}
	lines := make([]string, 0)
	for _, key := range keys {
		if n := counts[key]; n > 0 {
			lines = append(lines, fmt.Sprintf("%-8s %s %d", key, Bar(n, total, width), n))
		}
	}
	return strings.Join(lines, "\n")
}
func Percent(value, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(value) * 100 / float64(total)
}
func SummaryLine(label string, value int) string { return fmt.Sprintf("%s: %d", label, value) }

func RenderTable(rows []model.Record) string {
	lines := []string{"ID | 状态 | 标题"}
	for _, row := range rows {
		lines = append(lines, fmt.Sprintf("%s | %s | %s", row.ID, model.StatusLabel(row.Status), row.Title))
	}
	return strings.Join(lines, "\n")
}

func RenderCompact(record model.Record) string {
	return fmt.Sprintf("%s/%s/%s", record.ID, record.Status, record.Priority)
}

func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func NonEmptyLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}

func PadLabel(label string, width int) string {
	if len(label) >= width {
		return label
	}
	return label + strings.Repeat(" ", width-len(label))
}

func Clamp(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func Ratio(value, total int) string {
	return fmt.Sprintf("%.1f%%", Percent(value, total))
}

func EvenlySplit(total, parts int) []int {
	if parts <= 0 {
		return []int{}
	}
	result := make([]int, parts)
	for i := 0; i < total; i++ {
		result[i%parts]++
	}
	return result
}

func Sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func Average(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(Sum(values)) / float64(len(values))
}

func Zero() int { return 0 }
