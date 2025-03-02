package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zhengyansheng/workflow/config"
	"github.com/zhengyansheng/workflow/handler"
	"github.com/zhengyansheng/workflow/pkg/file"
)

var (
	// 这些变量会在编译时注入
	Version   string
	GitCommit string
	BuildTime string
)

// setupMiddleware 设置中间件
func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
}

// setupRoutes 设置路由
func setupRoutes(e *echo.Echo, workflowHandler *handler.WorkflowHandler) {
	// API路由
	v1 := e.Group("/api/v1")
	{
		workflows := v1.Group("/workflows")
		workflows.GET("", workflowHandler.ListWorkflows)
		workflows.POST("/:workflow_name/terminate", workflowHandler.TerminateWorkflow)
		workflows.GET("/:workflow_name/logs", workflowHandler.StreamLogs)
	}
}

// initServer 初始化服务器
func initServer(cfg *config.Config) (*echo.Echo, error) {
	// 读取kubeconfig
	kubeConfig, err := file.ReadLocalFile()
	if err != nil {
		return nil, err
	}

	// 创建Echo实例
	e := echo.New()

	// 设置中间件
	setupMiddleware(e)

	// 创建handler
	workflowHandler := handler.NewWorkflowHandler(kubeConfig, cfg.Workflow.Namespace)

	// 设置路由
	setupRoutes(e, workflowHandler)

	return e, nil
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	// 初始化服务器
	e, err := initServer(cfg)
	if err != nil {
		panic(fmt.Sprintf("初始化服务器失败: %v", err))
	}

	// 可以添加版本信息日志
	e.Logger.Printf("Starting workflow-server Version: %s, GitCommit: %s, BuildTime: %s\n",
		Version, GitCommit, BuildTime)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	e.Logger.Fatal(e.Start(serverAddr))
}
