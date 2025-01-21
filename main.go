package main

import (
	"fmt"
	"log"
	"os"

	"github.com/edwardOWO/goexample"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
)

func main() {

	utils.PrintMessage("Hello, Go!")
	result := utils.AddNumbers(5, 3)
	fmt.Println("Result of addition:", result)

	// 配置 kubeconfig 路徑
	kubeconfig := `/home/coder/goexample/test.config`
	// 检查 kubeconfig 文件是否存在
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Fatalf("kubeconfig 文件不存在: %v", err)
	}

	// 初始化 Helm CLI 设置
	settings := cli.New()
	settings.KubeConfig = kubeconfig

	// 初始化 Helm Action 配置
	actionConfig := new(action.Configuration)
	clientGetter := kube.GetConfig(settings.KubeConfig, "", settings.Namespace()) // 指定 default 命名空间
	err := actionConfig.Init(clientGetter, "", os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		log.Fatalf("无法初始化 Helm 配置: %v", err)
	}

	// 创建 List Client
	listClient := action.NewList(actionConfig)
	//listClient.Deployed = true // 仅查询已部署的 Releases
	listClient.AllNamespaces = true

	// 查询 Releases
	releases, err := listClient.Run()
	if err != nil {
		log.Fatalf("无法获取 Releases 列表: %v", err)
	}

	// 输出结果
	fmt.Println("Default 命名空间中的已部署 Releases:")
	if len(releases) == 0 {
		fmt.Println("没有找到已部署的 Release")
	} else {
		for _, r := range releases {
			fmt.Printf("- 名称: %s, 状态: %s\n", r.Name, r.Info.Status)
		}
	}
}
