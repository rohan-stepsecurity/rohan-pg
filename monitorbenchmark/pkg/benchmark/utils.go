package monitorbenchmark

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func (wr *WorkflowRunner) triggerWorkflow() error {
	args := []string{"workflow", "run", wr.config.WorkflowFile}
	// TODO!: Handle no harden runner runs
	if wr.config.WorkflowFile != "stress.yml" {
		args = append(args,
			"-f", "RunnerLabel="+wr.config.RunnerLabel,
			"-f", "UseHardenRunner="+wr.config.UseHardenRunner,
			"--repo", "harden-runner-canary/arc-tracer-benchmark",
		)
	}

	cmd := exec.Command("gh", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("workflow trigger failed: %w\n%s", err, string(output))
	}

	time.Sleep(3 * time.Second)
	return nil
}

func (wr *WorkflowRunner) getLatestRunID() (string, error) {
	cmd := exec.Command("gh", "run", "list",
		"--workflow="+wr.config.WorkflowFile,
		"--limit=1",
		"--json=databaseId",
		"--jq", ".[0].databaseId",
		"--repo", "harden-runner-canary/arc-tracer-benchmark",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get run ID: %w\n%s", err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

func (wr *WorkflowRunner) waitForWorkflow() error {
	cmd := exec.Command("gh", "run", "watch", *wr.runID, "--repo=harden-runner-canary/arc-tracer-benchmark")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func (wr *WorkflowRunner) getOutputPath(metricType string) string {
	switch metricType {
	case "node":
		switch wr.runID {
		case nil:
			return fmt.Sprintf("%s/%s-node-metrics.csv", wr.config.OutputDir, wr.fileID)
		default:
			return fmt.Sprintf("%s/%s-node-metrics.csv", wr.config.OutputDir, *wr.runID)
		}
	case "pod":
		switch wr.runID {
		case nil:
			return fmt.Sprintf("%s/%s-pod-metrics.csv", wr.config.OutputDir, wr.fileID)
		default:
			return fmt.Sprintf("%s/%s-pod-metrics.csv", wr.config.OutputDir, *wr.runID)
		}
	default:
		switch wr.runID {
		case nil:
			return fmt.Sprintf("%s/%s-details.json", wr.config.OutputDir, wr.fileID)
		default:
			return fmt.Sprintf("%s/%s-details.json", wr.config.OutputDir, *wr.runID)
		}

	}
}

func (wr *WorkflowRunner) renameFile(oldFilename, newFilename string) error {
	err := os.Rename(oldFilename, newFilename)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	return nil
}
