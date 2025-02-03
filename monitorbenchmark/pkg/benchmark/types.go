package monitorbenchmark

import (
	"time"

	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type Config struct {
	WorkflowFile    string
	RunnerLabel     string
	UseHardenRunner string
	OutputDir       string
	MetricsInterval time.Duration
	CooldownPeriod  time.Duration
	IsReport        bool
}

type WorkflowRunner struct {
	config        *Config
	metricsClient *versioned.Clientset
	runID         *string
	fileID        string
}

type Step struct {
	CompletedAt string `json:"completedAt"`
	Name        string `json:"name"`
	JobName     string `json:"jobName"`
}

type Job struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

type RunDetails struct {
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	CreatedAt    time.Time `json:"createdAt"`
	StartedAt    time.Time `json:"startedAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	DisplayTitle string    `json:"displayTitle"`
	Jobs         []Job     `json:"jobs"`
}
