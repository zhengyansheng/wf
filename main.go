package main

import (
	"time"

	"github.com/zhengyansheng/wf/logs"
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
	err := logs.WorkflowLogs(workflow, "", since, sinceTime, tailLines, grep, selector, container, follow, previous, timestamps, noColor)
	if err != nil {
		panic(err)
	}
}
