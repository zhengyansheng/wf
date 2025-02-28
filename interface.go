package wf

type WorkflowCommand interface {
	TerminateWorkflow(workflowName string) error
	ListWorkflows() error
	StreamLogs(workflow, podName string, tailLines int64) (logLineChan chan string, err error)
}
