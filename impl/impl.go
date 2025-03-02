package impl

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/zhengyansheng/workflow/api"
	"github.com/zhengyansheng/workflow/common"
	"github.com/zhengyansheng/workflow/list"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

type Client struct {
	ctx       context.Context
	client    apiclient.Client
	namespace string
}

func NewWorkflowClient(kubeConfig, namespace string) WorkflowCommand {
	// 初始化 API 客户端
	ctx := context.Background()
	ctx, apiClient, err := api.NewAPIClient(ctx, []byte(kubeConfig))
	if err != nil {
		return new(Client)
	}

	return &Client{
		ctx:       ctx,
		client:    apiClient,
		namespace: namespace,
	}
}

func (c *Client) TerminateWorkflow(workflowName string) error {

	serviceClient := c.client.NewWorkflowServiceClient()
	wf, err := serviceClient.TerminateWorkflow(c.ctx, &workflowpkg.WorkflowTerminateRequest{
		Name:      workflowName,
		Namespace: c.namespace,
	})
	if err != nil {
		return err
	}
	fmt.Printf("workflow %s terminated\n", wf.Name)
	return nil
}

func (c *Client) ListWorkflows() (wfv1.Workflows, error) {
	serviceClient := c.client.NewWorkflowServiceClient()
	workflows, err := list.ListWorkflows(c.ctx, serviceClient)
	if err != nil {
		return nil, err
	}
	return workflows, nil
}

func (c *Client) StreamLogs(workflow, podName string, tailLines int64) (chan string, error) {
	/*
		podName, tailLines 可留空即可
	*/

	logOptions := &corev1.PodLogOptions{
		Container: "main",
		Follow:    true,
	}

	if tailLines >= 0 {
		logOptions.TailLines = ptr.To(tailLines)
	}

	serviceClient := c.client.NewWorkflowServiceClient()
	return common.LogWorkflowWithChannel(c.ctx, serviceClient, c.namespace, workflow, podName, logOptions), nil
}
