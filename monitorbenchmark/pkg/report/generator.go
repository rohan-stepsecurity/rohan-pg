package report

import (
	"fmt"
	"strings"
	"time"
)

func GenerateReport(trials []Trial) string {
	var sb strings.Builder

	sb.WriteString("# ARC Harden Runner Benchmark Report\n\n")
	sb.WriteString(generateWorkflowSection(trials))
	sb.WriteString("\n---\n\n")
	sb.WriteString(generateResourceSection("Node", trials, func(t Trial) ResourceMetrics { return t.NodeMetrics }))
	sb.WriteString("\n---\n\n")
	sb.WriteString(generateResourceSection("Pod", trials, func(t Trial) ResourceMetrics { return t.PodMetrics }))
	sb.WriteString("\n---\n\n")
	//sb.WriteString(generateComparisonSection(trials))

	return sb.String()
}

func generateWorkflowSection(trials []Trial) string {
	var sb strings.Builder
	sb.WriteString("## Workflow Execution Times\n\n")

	// Group trials by runner type
	grouped := make(map[string][]Trial)
	for _, trial := range trials {
		grouped[trial.RunnerType] = append(grouped[trial.RunnerType], trial)
	}

	for runnerType, trials := range grouped {
		sb.WriteString(fmt.Sprintf("### %s\n", formatRunnerName(runnerType)))
		sb.WriteString("| Trial | Actual Run ID | Duration |\n")
		sb.WriteString("|-------|---------------|----------|\n")

		var totalDuration time.Duration
		for _, trial := range trials {
			duration := trial.WorkflowTime
			sb.WriteString(fmt.Sprintf("| %d | %s | %s |\n",
				trial.TrialNumber,
				trial.ActualRunID,
				formatDuration(duration),
			))
			totalDuration += duration
		}

		avg := totalDuration / time.Duration(len(trials))
		sb.WriteString(fmt.Sprintf("| **Avg** | | **%s** |\n\n", formatDuration(avg)))
	}

	return sb.String()
}

func generateResourceSection(resType string, trials []Trial, getMetrics func(Trial) ResourceMetrics) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s Resource Utilization\n\n", resType))

	grouped := make(map[string][]Trial)
	for _, trial := range trials {
		grouped[trial.RunnerType] = append(grouped[trial.RunnerType], trial)
	}

	for runnerType, trials := range grouped {
		sb.WriteString(fmt.Sprintf("### %s\n", formatRunnerName(runnerType)))
		sb.WriteString("| Trial | CPU Low (m) | CPU Avg (m) | CPU High (m) | Memory Low (Mi) | Memory Avg (Mi) | Memory High (Mi) |\n")
		sb.WriteString("|-------|-------------|-------------|--------------|-----------------|------------------|------------------|\n")

		for _, trial := range trials {
			metrics := getMetrics(trial)
			sb.WriteString(fmt.Sprintf("| %d | %.1f | %.1f | %.1f | %.1f | %.1f | %.1f |\n",
				trial.TrialNumber,
				metrics.CPU.Low,
				metrics.CPU.Avg,
				metrics.CPU.High,
				metrics.Memory.Low,
				metrics.Memory.Avg,
				metrics.Memory.High,
			))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func generateComparisonSection(trials []Trial) string {
	var sb strings.Builder
	sb.WriteString("## Performance Comparisons\n\n")

	sb.WriteString("### No Harden Runner vs Current Harden Runner\n")
	sb.WriteString("| Metric | No Harden | Current | Change |\n")
	sb.WriteString("|--------|-----------|---------|--------|\n")
	sb.WriteString("| Workflow Duration | 340s | 416s | +22.35% |\n")

	return sb.String()
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02dh%02dm%02ds", hours, minutes, seconds)
}

func formatRunnerName(runnerType string) string {
	return strings.Title(strings.ReplaceAll(runnerType, "-", " "))
}
