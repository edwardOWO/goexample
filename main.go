package main

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
)

func main() {
	// 配置 kubeconfig 路徑
	kubeconfig := `/home/coder/goexample/test.config`
	releaseName := "haproxy"                // 替換為你想查詢歷史記錄的 Release 名稱
	namespace := "default"       // 替換為 Release 所在的命名空間

	// 檢查 kubeconfig 文件是否存在
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Fatalf("kubeconfig 文件不存在: %v", kubeconfig)
	}

	// 初始化 Helm 的 CLI 設定
	settings := cli.New()
	settings.KubeConfig = kubeconfig
	settings.SetNamespace(namespace)

	// 初始化 Action Configuration
	actionConfig := new(action.Configuration)
	clientGetter := kube.GetConfig(settings.KubeConfig, "", settings.Namespace())
	err := actionConfig.Init(clientGetter, settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		log.Fatalf("無法初始化 Helm 設定: %v", err)
	}

	// 創建 History 動作
	historyClient := action.NewHistory(actionConfig)

	// 獲取 Release 的歷史記錄
	histories, err := historyClient.Run(releaseName)
	if err != nil {
		log.Fatalf("無法取得 Release 的歷史記錄: %v", err)
	}

	// 輸出歷史記錄
	fmt.Printf("Release '%s' 的歷史記錄:\n", releaseName)
	if len(histories) == 0 {
		fmt.Println("無歷史記錄")
	} else {
		for _, history := range histories {
			fmt.Printf("- 修訂版本: %d, 狀態: %s, 更新時間: %s\n",
				history.Version,
				history.Info.Status,
				history.Info.LastDeployed,
			)
		}
	}
}
