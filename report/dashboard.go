package report

import (
	"fmt"
	"labops/model"
	"labops/workflow"
	"sort"
	"strings"
	"time"
)

type Dashboard struct {
	GeneratedAt  time.Time
	Metrics      workflow.Metrics
	Attention    []model.Record
	TopAssignees []string
}

func BuildDashboard(records []model.Record, now time.Time) Dashboard {
	return Dashboard{GeneratedAt: now, Metrics: workflow.BuildMetrics(records, now), Attention: workflow.AttentionList(records), TopAssignees: workflow.TopAssignees(records)}
}
func RenderDashboard(d Dashboard) string {
	var b strings.Builder
	fmt.Fprintf(&b, "实验设备运维看板\n生成时间: %s\n", d.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "总数: %d 开放: %d 已关闭: %d\n", d.Metrics.Total, d.Metrics.Open, d.Metrics.Closed)
	b.WriteString("状态分布:\n")
	for _, key := range sortedKeys(d.Metrics.ByStatus) {
		fmt.Fprintf(&b, "- %s: %d\n", model.StatusLabel(key), d.Metrics.ByStatus[key])
	}
	b.WriteString("优先级分布:\n")
	for _, key := range sortedKeys(d.Metrics.ByPriority) {
		fmt.Fprintf(&b, "- %s: %d\n", key, d.Metrics.ByPriority[key])
	}
	fmt.Fprintf(&b, "需关注: %d\n", len(d.Attention))
	if len(d.TopAssignees) > 0 {
		fmt.Fprintf(&b, "负责人: %s\n", strings.Join(d.TopAssignees, ", "))
	}
	return b.String()
}
func sortedKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
func GroupByEquipment(records []model.Record) map[string][]model.Record {
	groups := make(map[string][]model.Record)
	for _, r := range records {
		groups[r.EquipmentID] = append(groups[r.EquipmentID], r)
	}
	return groups
}
func GroupByAssignee(records []model.Record) map[string][]model.Record {
	groups := make(map[string][]model.Record)
	for _, r := range records {
		key := r.Assignee
		if key == "" {
			key = "unassigned"
		}
		groups[key] = append(groups[key], r)
	}
	return groups
}
func AgingBuckets(records []model.Record, now time.Time) map[string]int {
	buckets := map[string]int{"0-1d": 0, "2-7d": 0, "8-30d": 0, "31d+": 0}
	for _, r := range records {
		age := now.Sub(r.UpdatedAt)
		if age <= 24*time.Hour {
			buckets["0-1d"]++
		} else if age <= 7*24*time.Hour {
			buckets["2-7d"]++
		} else if age <= 30*24*time.Hour {
			buckets["8-30d"]++
		} else {
			buckets["31d+"]++
		}
	}
	return buckets
}
func FormatAging(buckets map[string]int) string {
	keys := []string{"0-1d", "2-7d", "8-30d", "31d+"}
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, buckets[k]))
	}
	return strings.Join(parts, " ")
}
func StatusLine(status string, count int) string {
	return fmt.Sprintf("%-12s %4d", model.StatusLabel(status), count)
}
func SortRecords(records []model.Record) []model.Record { return model.SortByUpdated(records, true) }
func LimitAttention(records []model.Record, limit int) []model.Record {
	if limit < 1 {
		return []model.Record{}
	}
	if limit >= len(records) {
		return records
	}
	return records[:limit]
}
