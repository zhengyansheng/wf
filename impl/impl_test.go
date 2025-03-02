package impl

import (
	"fmt"
	"github.com/zhengyansheng/workflow/pkg/file"
	"testing"
)

func TestImpl(t *testing.T) {
	kubeConfig, err := file.ReadLocalFile()
	if err != nil {
		t.Error(err)
	}

	w := NewWorkflowClient(kubeConfig, "argo")
	//w.ListWorkflows()

	logLineChan, err := w.StreamLogs("java-pipeline-glpln", "", int64(-1))
	if err != nil {
		t.Failed()
	}

	for {
		line, ok := <-logLineChan
		if !ok {
			break
		}
		fmt.Println(line)
	}

}
