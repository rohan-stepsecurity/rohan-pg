package monitorbenchmark

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type CLIOptions struct {
	WorkflowFile    string
	RunnerLabel     string
	UseHardenRunner string
	OutputDir       string
	MetricsInterval time.Duration
	Cooldown        time.Duration
	IsReport        bool
}

func ParseFlags() *CLIOptions {
	var opts CLIOptions

	flag.StringVar(&opts.WorkflowFile, "workflow-file", "", "Workflow file name (required)")
	flag.StringVar(&opts.RunnerLabel, "runner-label", "", "Runner label (required)")
	flag.StringVar(&opts.UseHardenRunner, "use-harden-runner", "true", "Use hardened runner")
	flag.StringVar(&opts.OutputDir, "output-dir", "./metrics", "Output directory")
	flag.DurationVar(&opts.MetricsInterval, "metrics-interval", 2*time.Second, "Metrics collection interval")
	flag.DurationVar(&opts.Cooldown, "cooldown", 30*time.Second, "Cooldown period after workflow completion")
	flag.BoolVar(&opts.IsReport, "is-report", false, "To generate report or not")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Println(opts)

	return &opts
}

func CLIOptionsToConfig(opts *CLIOptions) *Config {
	cfg := &Config{
		WorkflowFile:    opts.WorkflowFile,
		RunnerLabel:     opts.RunnerLabel,
		UseHardenRunner: opts.UseHardenRunner,
		OutputDir:       opts.OutputDir,
		MetricsInterval: opts.MetricsInterval,
		CooldownPeriod:  opts.Cooldown,
		IsReport:        opts.IsReport,
	}
	return cfg

}

func ValidateConfig(cfg *Config) error {
	if cfg.IsReport {
		return nil
	}
	if cfg.WorkflowFile == "" {
		return fmt.Errorf("workflow-file is required")
	}
	if cfg.RunnerLabel == "" {
		return fmt.Errorf("runner-label is required")
	}

	return nil
}
