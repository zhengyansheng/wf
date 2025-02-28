package terminate

import (
	"context"
	"fmt"
	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/zhengyansheng/wf/api"
)

type workflowLogsClient struct {
	KubeConfig string `json:"kube_config"`
	Namespace  string `json:"namespace"`
}

func NewWorkflowLogsClient(kubeConfig, namespace string) *workflowLogsClient {
	return &workflowLogsClient{kubeConfig, namespace}
}

func (w *workflowLogsClient) TerminateWorkflow(workflowName string) error {
	// 初始化 API 客户端
	ctx := context.Background()
	ctx, apiClient, err := api.NewAPIClient(ctx, []byte(w.KubeConfig))
	if err != nil {
		return err
	}

	serviceClient := apiClient.NewWorkflowServiceClient()
	wf, err := serviceClient.TerminateWorkflow(ctx, &workflowpkg.WorkflowTerminateRequest{
		Name:      workflowName,
		Namespace: w.Namespace,
	})
	if err != nil {
		return err
	}
	fmt.Printf("workflow %s terminated\n", wf.Name)
	return nil
}
