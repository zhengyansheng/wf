package impl

import (
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

type WorkflowCommand interface {
	TerminateWorkflow(workflowName string) error
	ListWorkflows() (wfv1.Workflows, error)
	StreamLogs(workflow, podName string, tailLines int64) (logLineChan chan string, err error)
}
