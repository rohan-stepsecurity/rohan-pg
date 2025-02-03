package monitorbenchmark

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func (wr *WorkflowRunner) saveArtifacts(startTime time.Time) error {
	fmt.Println(startTime)
	detailsPath := wr.getOutputPath("json")
	if err := wr.saveRunDetails(detailsPath); err != nil {
		return err
	}

	steps, err := wr.readSteps(detailsPath)
	if err != nil {
		return err
	}

	if err := wr.processMetricsFiles(steps); err != nil {
		return err
	}

	return nil
}

func (wr *WorkflowRunner) saveRunDetails(path string) error {
	cmd := exec.Command("gh", "run", "view", *wr.runID,
		"--json", "status,conclusion,createdAt,displayTitle,jobs,startedAt,updatedAt",
		"--repo=harden-runner-canary/arc-tracer-benchmark",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch workflow details: %w\n%s", err, string(output))
	}

	return os.WriteFile(path, output, 0644)
}

func (wr *WorkflowRunner) readSteps(path string) ([]Step, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var details struct{ Jobs []Job }
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, err
	}

	var steps []Step
	for _, job := range details.Jobs {
		for _, step := range job.Steps {
			step.JobName = job.Name
			steps = append(steps, step)
		}
	}

	sort.Slice(steps, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, steps[i].CompletedAt)
		t2, _ := time.Parse(time.RFC3339, steps[j].CompletedAt)
		return t1.Before(t2)
	})

	return steps, nil
}

func (wr *WorkflowRunner) processMetricsFiles(steps []Step) error {
	patterns := []string{"*-node-metrics.csv", "*-pod-metrics.csv"}

	for _, pattern := range patterns {
		files, err := filepath.Glob(filepath.Join(wr.config.OutputDir, pattern))
		if err != nil {
			return err
		}
		fmt.Println("files found to update, ", files)

		for _, file := range files {
			fmt.Printf("Updating file %s\n", file)
			if err := wr.updateCSVWithSteps(file, steps); err != nil {
				fmt.Println("error while updating file", err)
				return err
			}
			fmt.Printf("renaming file %s\n", file)
			if err := wr.renameFile(file, strings.ReplaceAll(file, wr.fileID, *wr.runID)); err != nil {
				fmt.Println("error while renaming file", err)
				return err
			}
		}
	}

	return nil
}

func (wr *WorkflowRunner) updateCSVWithSteps(path string, steps []Step) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	file.Close()

	for i, record := range records {
		if i == 0 {
			continue
		}

		ts, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			continue
		}

		for j := 0; j < len(steps); j++ {
			stepTime, _ := time.Parse(time.RFC3339, steps[j].CompletedAt)
			if ts.Before(stepTime) {
				record[2] = steps[j].JobName
				record[3] = steps[j].Name
				break
			}
		}
	}

	outputFile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func (wr *WorkflowRunner) cleanup() error {
	patterns := []string{
		fmt.Sprintf("%s/*-node-metrics.csv", wr.config.OutputDir),
		fmt.Sprintf("%s/*-pod-metrics.csv", wr.config.OutputDir),
		fmt.Sprintf("%s/*-details.json", wr.config.OutputDir),
	}

	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, file := range files {
			if strings.Contains(file, wr.fileID) {
				os.Remove(file)
			}
		}
	}
	return nil
}
