package logs

import (
	"context"
	"errors"
	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	"github.com/zhengyansheng/workflow/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"time"
)

func SteamLogs(workflow, podName string, since time.Duration, sinceTime string, tailLines int64, grep, selector, container string, follow, previous, timestamps, noColor bool) error {
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

	return common.LogWorkflowWithColor(ctx, serviceClient, namespace, workflow, podName, grep, selector, logOptions)
	//return common.LogWorkflowWithColor(ctx, serviceClient, namespace, logs, podName, grep, selector, logOptions)
}
