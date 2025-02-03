package main

import (
	"log"
	monitorbenchmark "monitorbenchmark/pkg/benchmark"
	"monitorbenchmark/pkg/report"
	"os"
)

func main() {
	cliOpts := monitorbenchmark.ParseFlags()

	cfg := monitorbenchmark.CLIOptionsToConfig(cliOpts)

	if err := monitorbenchmark.ValidateConfig(cfg); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	if cfg.IsReport {

		trials, err := report.ParseData("../artifact")
		if err != nil {
			log.Fatalf("Error parsing data: %v", err)
		}

		reportContent := report.GenerateReport(trials)

		if err := os.WriteFile("REPORT.md", []byte(reportContent), 0644); err != nil {
			log.Fatalf("Error writing report: %v", err)
		}

		return

	}

	runner, err := monitorbenchmark.NewWorkflowRunner(cfg)
	if err != nil {
		log.Fatalf("Failed to create workflow runner: %v", err)
	}

	if err := runner.Execute(); err != nil {
		log.Fatalf("Workflow execution failed: %v", err)
	}
}
