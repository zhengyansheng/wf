package logs

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkflowLogs(t *testing.T) {
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
	err := WorkflowLogs(workflow, "", since, sinceTime, tailLines, grep, selector, container, follow, previous, timestamps, noColor)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestWorkflowLogsWithChannel(t *testing.T) {
	// 设置默认值
	tailLines := int64(-1)
	workflow := "java-pipeline-67lsm"

	// 调用日志函数
	logLineChan, err := WorkflowLogsWithChannel(workflow, "", tailLines)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for {
		line, ok := <-logLineChan
		if !ok {
			break
		}
		fmt.Println(line)
	}

}
