package monitorbenchmark

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	defaultMetricsInterval = 2 * time.Second
	defaultCooldown        = 30 * time.Second
)

func NewWorkflowRunner(cfg *Config) (*WorkflowRunner, error) {
	if cfg.MetricsInterval == 0 {
		cfg.MetricsInterval = defaultMetricsInterval
	}
	if cfg.CooldownPeriod == 0 {
		cfg.CooldownPeriod = defaultCooldown
	}

	token, ok := os.LookupEnv("KUBE_TOKEN")
	if !ok {
		slog.Error("Failed to get KUBE_TOKEN ")

	}

	host, err := getKubeHost()
	if err != nil {
		slog.Error("Failed to get KUBE_HOST ", slog.Any("err", err))
		return nil, errors.New(`failed to get KUBE_HOST`)
	}

	config := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	return &WorkflowRunner{
		config:        cfg,
		metricsClient: metricsClient,
		fileID:        generateRandomString(6),
	}, nil
}

func getKubeHost() (string, error) {
	// Load the kubeconfig from the default location (~/.kube/config)
	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Get the current context
	currentContext := config.CurrentContext
	if currentContext == "" {
		return "", fmt.Errorf("no current context found in kubeconfig")
	}

	// Get the cluster associated with the current context
	context, exists := config.Contexts[currentContext]
	if !exists {
		return "", fmt.Errorf("context %q not found in kubeconfig", currentContext)
	}

	cluster, exists := config.Clusters[context.Cluster]
	if !exists {
		return "", fmt.Errorf("cluster %q not found in kubeconfig", context.Cluster)
	}

	// Return the Kubernetes API server URL
	return cluster.Server, nil
}

func (wr *WorkflowRunner) Execute() error {
	err := os.Mkdir("metrics", 0755)
	fmt.Println(err)
	err = os.Mkdir("./metrics", 0755)
	fmt.Println(err)

	start := time.Now()
	defer wr.cleanup()

	if err := wr.triggerWorkflow(); err != nil {
		return err
	}

	runID, err := wr.getLatestRunID()
	if err != nil {
		return err
	}
	wr.runID = &runID

	nodeDone := make(chan struct{})
	podDone := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	go wr.monitorNodeMetrics(nodeDone, &wg)
	go wr.monitorPodMetrics(podDone, &wg)

	if err := wr.waitForWorkflow(); err != nil {
		return err
	}

	fmt.Println("waiting for workflow")

	close(nodeDone)
	close(podDone)
	wg.Wait()

	fmt.Println("waiting for workflow done, closed channels")

	return wr.saveArtifacts(start)
}
