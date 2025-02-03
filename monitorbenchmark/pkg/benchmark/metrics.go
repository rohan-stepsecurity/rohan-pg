package monitorbenchmark

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wr *WorkflowRunner) monitorPodMetrics(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	filePath := wr.getOutputPath("pod")
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"PodName", "Timestamp", "CurrentJob", "CurrentStep", "Container", "CPU(m)", "Memory(Mi)"})

	for {
		select {
		case <-done:
			return
		default:
			metrics, err := wr.metricsClient.MetricsV1beta1().PodMetricses("kube-system").List(
				context.Background(),
				metav1.ListOptions{LabelSelector: "app.kubernetes.io/instance=hardenrunner"},
			)

			if err == nil {
				timestamp := time.Now().UTC().Format(time.RFC3339)
				for _, pod := range metrics.Items {
					for _, container := range pod.Containers {
						cpu := container.Usage.Cpu().MilliValue()
						memoryMi := float64(container.Usage.Memory().Value()) / (1024 * 1024)
						writer.Write([]string{
							pod.GetName(),
							timestamp,
							"NA",
							"NA",
							container.Name,
							fmt.Sprintf("%d", cpu),
							fmt.Sprintf("%.3f", memoryMi),
						})
					}
				}
				writer.Flush()
			}
			time.Sleep(wr.config.MetricsInterval)
		}
	}
}

func (wr *WorkflowRunner) monitorNodeMetrics(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	filePath := wr.getOutputPath("node")
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ClusterVersion", "Timestamp", "CurrentJob", "CurrentStep", "NodeName", "CPU(m)", "Memory(Mi)"})

	for {
		select {
		case <-done:
			return
		default:
			metrics, err := wr.metricsClient.MetricsV1beta1().NodeMetricses().List(
				context.Background(),
				metav1.ListOptions{},
			)

			if err == nil {
				timestamp := time.Now().UTC().Format(time.RFC3339)
				version, _ := wr.metricsClient.ServerVersion()

				for _, node := range metrics.Items {
					cpu := node.Usage.Cpu().MilliValue()
					memoryMi := float64(node.Usage.Memory().Value()) / (1024 * 1024)

					writer.Write([]string{
						version.GitVersion,
						timestamp,
						"NA",
						"NA",
						node.GetName(),
						fmt.Sprintf("%d", cpu),
						fmt.Sprintf("%.3f", memoryMi),
					})
				}
				writer.Flush()
			}
			time.Sleep(wr.config.MetricsInterval)
		}
	}
}
