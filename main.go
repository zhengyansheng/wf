package main

import (
	"context"
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	"github.com/zhengyansheng/wf/common"
)

func main() {
	// 设置默认值
	since := time.Duration(0)
	sinceTime := ""
	tailLines := int64(-1)
	grep := ""
	selector := ""
	container := "main"
	follow := true
	previous := false
	timestamps := false
	noColor := false
	workflow := "java-pipeline-67lsm"

	// 调用日志函数
	err := getWorkflowLogs(workflow, "", since, sinceTime, tailLines, grep, selector, container, follow, previous, timestamps, noColor)
	if err != nil {
		panic(err)
	}
}

func getWorkflowLogs(workflow, podName string, since time.Duration, sinceTime string, tailLines int64, grep, selector, container string, follow, previous, timestamps, noColor bool) error {
	logOptions := &corev1.PodLogOptions{
		Container:  container,
		Follow:     follow,
		Previous:   previous,
		Timestamps: timestamps,
	}

	if since > 0 && sinceTime != "" {
		return errors.New("--since-time and --since cannot be used together")
	}

	if since > 0 {
		logOptions.SinceSeconds = ptr.To(int64(since.Seconds()))
	}

	if sinceTime != "" {
		parsedTime, err := time.Parse(time.RFC3339, sinceTime)
		if err != nil {
			return err
		}
		sinceTime := metav1.NewTime(parsedTime)
		logOptions.SinceTime = &sinceTime
	}

	if tailLines >= 0 {
		logOptions.TailLines = ptr.To(tailLines)
	}

	// 初始化 API 客户端
	ctx := context.Background()
	ctx, apiClient, err := client.NewAPIClient(ctx)
	if err != nil {
		return err
	}
	serviceClient := apiClient.NewWorkflowServiceClient()
	namespace := "argo" // 设置默认命名空间

	common.NoColor = noColor

	return common.LogWorkflow(ctx, serviceClient, namespace, workflow, podName, grep, selector, logOptions)
	//return common.LogWorkflow(ctx, serviceClient, namespace, workflow, podName, grep, selector, logOptions)
}
