package common

import (
	"context"
	"fmt"
	"io"
	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"

	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
)

func LogWorkflowWithColor(ctx context.Context, serviceClient workflowpkg.WorkflowServiceClient, namespace, workflow, podName, grep, selector string, logOptions *corev1.PodLogOptions) error {
	// logs
	stream, err := serviceClient.WorkflowLogs(ctx, &workflowpkg.WorkflowLogRequest{
		Name:       workflow,
		Namespace:  namespace,
		PodName:    podName,
		LogOptions: logOptions,
		Selector:   selector,
		Grep:       grep,
	})
	if err != nil {
		return err
	}

	// loop on log lines
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(ansiFormat(fmt.Sprintf("%s: %s", event.PodName, event.Content), ansiColorCode(event.PodName)))
	}
}

func LogWorkflowWithChannel(ctx context.Context, serviceClient workflowpkg.WorkflowServiceClient, namespace, workflow, podName string, logOptions *corev1.PodLogOptions) (logLineChan chan string) {
	logLineChan = make(chan string)

	go func() {
		defer close(logLineChan)
		// logs
		stream, err := serviceClient.WorkflowLogs(ctx, &workflowpkg.WorkflowLogRequest{
			Name:       workflow,
			Namespace:  namespace,
			PodName:    podName,
			LogOptions: logOptions,
		})
		if err != nil {
			klog.Errorf("failed to get stream logs: %v", err)
			return
		}

		// loop on log lines
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				klog.Errorf("failed to get logs logs: %v", err)
				return
			}
			logLineChan <- ansiFormat(event.Content)
		}
	}()

	return logLineChan
}
