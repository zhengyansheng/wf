package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zhengyansheng/workflow/impl"
)

type WorkflowHandler struct {
	client impl.WorkflowCommand
}

func NewWorkflowHandler(kubeConfig, namespace string) *WorkflowHandler {
	return &WorkflowHandler{
		client: impl.NewWorkflowClient(kubeConfig, namespace),
	}
}

// ListWorkflows godoc
// @Summary 列出所有工作流
// @Description 获取所有工作流列表
// @Tags workflow
// @Accept json
// @Produce json
// @Success 200 {array} wfv1.Workflow
// @Router /api/v1/workflows [get]
func (h *WorkflowHandler) ListWorkflows(c echo.Context) error {
	workflows, err := h.client.ListWorkflows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// 构造响应数据
	response := make([]map[string]interface{}, 0, len(workflows))
	for _, wf := range workflows {
		// 计算工作流运行时长
		var duration string
		if !wf.Status.StartedAt.IsZero() {
			if !wf.Status.FinishedAt.IsZero() {
				duration = wf.Status.FinishedAt.Time.Sub(wf.Status.StartedAt.Time).String()
			}
		}

		// 计算工作流年龄
		var age string
		if !wf.CreationTimestamp.IsZero() {
			age = time.Since(wf.CreationTimestamp.Time).Round(time.Second).String()
		}

		workflow := map[string]interface{}{
			"name":       wf.Name,
			"namespace":  wf.Namespace,
			"status":     wf.Status.Phase,
			"age":        age,
			"duration":   duration,
			"priority":   wf.Spec.Priority,
			"message":    wf.Status.Message,
			"startedAt":  wf.Status.StartedAt,
			"finishedAt": wf.Status.FinishedAt,
		}
		response = append(response, workflow)
	}

	return c.JSON(http.StatusOK, response)
}

// TerminateWorkflow godoc
// @Summary 终止工作流
// @Description 终止指定的工作流
// @Tags workflow
// @Accept json
// @Produce json
// @Param workflow_name path string true "工作流名称"
// @Success 200 {string} string "成功"
// @Router /api/v1/workflows/{workflow_name}/terminate [post]
func (h *WorkflowHandler) TerminateWorkflow(c echo.Context) error {
	workflowName := c.Param("workflow_name")
	err := h.client.TerminateWorkflow(workflowName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "workflow terminated successfully",
	})
}

// StreamLogs godoc
// @Summary 流式获取工作流日志
// @Description 获取指定工作流的实时日志
// @Tags workflow
// @Accept json
// @Produce json
// @Param workflow_name path string true "工作流名称"
// @Param pod_name query string false "Pod名称"
// @Param tail_lines query int false "返回的日志行数"
// @Success 200 {string} string "成功"
// @Router /api/v1/workflows/{workflow_name}/logs [get]
func (h *WorkflowHandler) StreamLogs(c echo.Context) error {
	workflowName := c.Param("workflow_name")
	podName := c.QueryParam("pod_name")
	tailLines := int64(-1)

	if tl := c.QueryParam("tail_lines"); tl != "" {
		if v, err := strconv.ParseInt(tl, 10, 64); err == nil {
			tailLines = v
		}
	}

	logChan, err := h.client.StreamLogs(workflowName, podName, tailLines)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	for log := range logChan {
		if _, err := c.Response().Write([]byte("data: " + log + "\n\n")); err != nil {
			return err
		}
		c.Response().Flush()
	}

	return nil
}
