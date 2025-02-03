package report

import "time"

type Trial struct {
	ParentRunID  string
	ActualRunID  string
	TrialNumber  int
	RunnerType   string
	WorkflowTime time.Duration
	NodeMetrics  ResourceMetrics
	PodMetrics   ResourceMetrics
}

type ResourceMetrics struct {
	CPU struct {
		Low  float64
		Avg  float64
		High float64
	}
	Memory struct {
		Low  float64
		Avg  float64
		High float64
	}
}

type Comparison struct {
	BaselineRunner string
	TargetRunner   string
	Metrics        []MetricComparison
}

type MetricComparison struct {
	Name         string
	Baseline     float64
	Target       float64
	PercentDelta float64
}

type WorkflowDetails struct {
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Jobs       []Job     `json:"jobs"`
}

type Step struct {
	CompletedAt time.Time `json:"completedAt"`
	Name        string    `json:"name"`
	JobName     string    `json:"jobName"`
}

type Job struct {
	Name        string    `json:"name"`
	Steps       []Step    `json:"steps"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
}
