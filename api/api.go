package api

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

var (
	instanceID string
	//explicitPath string = "/Users/mac/.kube/xt-test-config"
	explicitPath string
)

func NewAPIClient(ctx context.Context, kubeConfig []byte) (context.Context, apiclient.Client, error) {
	clientConfig, err := getClientConfig(kubeConfig)
	if err != nil {
		return nil, nil, err
	}
	return apiclient.NewClientFromOpts(
		apiclient.Opts{
			ArgoServerOpts: client.ArgoServerOpts,
			InstanceID:     instanceID,
			AuthSupplier: func() string {
				authString, err := client.GetAuthString()
				if err != nil {
					log.Fatal(err)
				}
				return authString
			},
			ClientConfigSupplier: func() clientcmd.ClientConfig {
				return clientConfig
			},
			Offline:      client.Offline,
			OfflineFiles: client.OfflineFiles,
			Context:      ctx,
		})
}

func getClientConfig(kubeConfig []byte) (clientcmd.ClientConfig, error) {
	overridingClientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	return overridingClientConfig, nil
}

func GetConfig2() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	loadingRules.ExplicitPath = explicitPath
	return clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}, os.Stdin)
}

// CreateHiddenTempFile 在家目录下创建一个隐藏的临时文件，并将内容写入文件
// 返回文件路径和可能的错误
func CreateHiddenTempFile(content string) (string, error) {
	// 获取当前用户的家目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	// 在家目录下创建一个隐藏的临时文件
	tmpFile, err := ioutil.TempFile(homeDir, ".kubeconfig-config")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	// 将内容写入临时文件
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		os.Remove(tmpFile.Name()) // 如果写入失败，清理临时文件
		return "", fmt.Errorf("failed to write content to temporary file: %v", err)
	}

	// 返回临时文件的路径
	return fmt.Sprintf("%v/%v", homeDir, tmpFile.Name()), nil
}
