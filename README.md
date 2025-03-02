# argo workflow 

## API

```bash
curl http://localhost:8080/api/v1/workflows | jq .
curl http://localhost:8080/api/v1/workflows/java-pipeline-glpln/logs\?tail_lines\=100

```

## 构建

```bash
# 构建指定版本
make docker-build VERSION=v0.0.1

# 运行容器
make docker-run VERSION=v0.0.1
```




## argo logs -f <workflowName>

> 参考

https://github.com/argoproj/argo-workflows/blob/main/cmd/argo/commands/logs.go
https://github.com/argoproj/argo-workflows/blob/main/cmd/argo/commands/common/logs.go

## Client Library

https://argo-workflows.readthedocs.io/en/latest/client-libraries/
https://github.com/argoproj/argo-workflows/blob/main/examples/example-golang/main.go