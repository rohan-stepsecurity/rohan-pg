
BINARY_NAME := app
BUILD_DIR   := bin

# Environment variables with defaults
WORKFLOW_FILE ?= stress-harden-runner.yml
RUNNER_LABEL ?= no-harden-runner
OUTPUT_DIR ?= ./metrics
METRICS_INTERVAL ?= 2s
COOLDOWN ?= 30s
USE_HARDEN_RUNNER ?= false
IS_REPORT ?= false

.PHONY: build
build: 
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./app/

.PHONY: run
run: build
	@echo "Starting workflow runner with configuration:"
	@echo "Workflow File: $(WORKFLOW_FILE)"
	@echo "Runner Label: $(RUNNER_LABEL)"
	@echo "Output Directory: $(OUTPUT_DIR)"
	@echo "Metrics Interval: $(METRICS_INTERVAL)"
	@echo "Cooldown Period: $(COOLDOWN)"
	@$(BUILD_DIR)/$(BINARY_NAME) \
		--workflow-file $(WORKFLOW_FILE) \
		--runner-label $(RUNNER_LABEL) \
		--use-harden-runner $(USE_HARDEN_RUNNER) \
		--output-dir $(OUTPUT_DIR) \
		--metrics-interval $(METRICS_INTERVAL) \
		--cooldown $(COOLDOWN)


.PHONY: report
report: build
	@echo "Starting workflow runner with configuration:"
	@echo "Workflow File: $(WORKFLOW_FILE)"
	@echo "Runner Label: $(RUNNER_LABEL)"
	@echo "Output Directory: $(OUTPUT_DIR)"
	@echo "Metrics Interval: $(METRICS_INTERVAL)"
	@echo "Cooldown Period: $(COOLDOWN)"
	@echo "Is Report: $(IS_REPORT)"
	@$(BUILD_DIR)/$(BINARY_NAME) \
		--workflow-file $(WORKFLOW_FILE) \
		--runner-label $(RUNNER_LABEL) \
		--use-harden-runner $(USE_HARDEN_RUNNER) \
		--output-dir $(OUTPUT_DIR) \
		--metrics-interval $(METRICS_INTERVAL) \
		--cooldown $(COOLDOWN) \
		--is-report true