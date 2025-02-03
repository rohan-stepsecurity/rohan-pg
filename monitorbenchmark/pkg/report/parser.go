package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseData(rootDir string) ([]Trial, error) {
	var trials []Trial

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, fmt.Errorf("error reading root directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		parts := strings.Split(dirName, "-")
		if len(parts) < 6 {
			continue // Skip invalid directories
		}

		// Parse directory name: workflow-metrics-<runner-type>-<parent-run-id>-<trial-number>
		runnerType := strings.Join(parts[2:len(parts)-2], "-")
		parentRunID := parts[len(parts)-2]
		trialNumber, _ := strconv.Atoi(parts[len(parts)-1])

		trialDir := filepath.Join(rootDir, dirName)
		trial, err := parseTrialDir(trialDir, runnerType, parentRunID, trialNumber)
		if err != nil {
			return nil, fmt.Errorf("error parsing trial directory %s: %w", dirName, err)
		}

		trials = append(trials, *trial)
	}

	return trials, nil
}

func parseTrialDir(dirPath, runnerType, parentRunID string, trialNumber int) (*Trial, error) {
	var trial Trial
	trial.RunnerType = runnerType
	trial.ParentRunID = parentRunID
	trial.TrialNumber = trialNumber

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "-details.json") {
			actualRunID := strings.TrimSuffix(file.Name(), "-details.json")
			trial.ActualRunID = actualRunID

			// Parse workflow details
			detailsPath := filepath.Join(dirPath, file.Name())
			details, err := parseWorkflowDetails(detailsPath)
			if err != nil {
				return nil, err
			}
			for index := range details.Jobs {
				jobDuration := details.Jobs[index].CompletedAt.Sub(details.Jobs[index].StartedAt)
				trial.WorkflowTime += jobDuration
			}
		}

		if strings.HasSuffix(file.Name(), "-node-metrics.csv") {
			//actualRunID := strings.TrimSuffix(file.Name(), "-node-metrics.csv")
			metricsPath := filepath.Join(dirPath, file.Name())
			nodeMetrics, err := parseNodeMetrics(metricsPath)
			if err != nil {
				return nil, err
			}
			trial.NodeMetrics = *nodeMetrics
		}

		if strings.HasSuffix(file.Name(), "-pod-metrics.csv") {
			//actualRunID := strings.TrimSuffix(file.Name(), "-pod-metrics.csv")
			metricsPath := filepath.Join(dirPath, file.Name())
			podMetrics, err := parsePodMetrics(metricsPath)
			if err != nil {
				return nil, err
			}
			trial.PodMetrics = *podMetrics
		}
	}

	return &trial, nil
}

func parseWorkflowDetails(path string) (*WorkflowDetails, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var details WorkflowDetails
	if err := json.NewDecoder(file).Decode(&details); err != nil {
		return nil, err
	}
	return &details, nil
}

func parseNodeMetrics(path string) (*ResourceMetrics, error) {
	return parseMetricsCSV(path, []string{"ClusterVersion", "Timestamp", "CurrentJob", "CurrentStep", "NodeName", "CPU(m)", "Memory(Mi)"}, "node")
}

func parsePodMetrics(path string) (*ResourceMetrics, error) {
	return parseMetricsCSV(path, []string{"PodName", "Timestamp", "CurrentJob", "CurrentStep", "Container", "CPU(m)", "Memory(Mi)"}, "pod")
}

func parseMetricsCSV(path string, headers []string, fileType string) (*ResourceMetrics, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Validate headers
	if !validateHeaders(records[0], headers) {
		return nil, fmt.Errorf("invalid CSV headers in %s", path)
	}

	var cpuValues []float64
	var memValues []float64

	for i, record := range records[1:] {
		var cpuIndex, memIndex int
		if fileType == "node" {
			cpuIndex, memIndex = 5, 6
		} else if fileType == "pod" {
			cpuIndex, memIndex = 6, 7

		}

		cpu, err := strconv.ParseFloat(record[cpuIndex], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid CPU value in row %d: %w", i+1, err)
		}

		mem, err := strconv.ParseFloat(record[memIndex], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Memory value in row %d: %w", i+1, err)
		}

		cpuValues = append(cpuValues, cpu)
		memValues = append(memValues, mem)
	}

	return &ResourceMetrics{
		CPU:    calculateStats(cpuValues),
		Memory: calculateStats(memValues),
	}, nil
}

func validateHeaders(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i] != expected[i] {
			return false
		}
	}
	return true
}

func calculateStats(values []float64) struct {
	Low  float64
	Avg  float64
	High float64
} {
	if len(values) == 0 {
		return struct{ Low, Avg, High float64 }{}
	}

	low := values[0]
	high := values[0]
	sum := 0.0

	for _, v := range values {
		if v < low {
			low = v
		}
		if v > high {
			high = v
		}
		sum += v
	}

	return struct{ Low, Avg, High float64 }{
		Low:  low,
		High: high,
		Avg:  sum / float64(len(values)),
	}
}
