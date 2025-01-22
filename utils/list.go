package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Release 定义一个结构体用于存储 Helm Release 的信息
type Release struct {
	Name       string `json:"name"`       // Release 名称
	Status     string `json:"status"`     // Release 状态
	AppVersion string `json:"appversion"` // Release 状态
	Version    int    `json:"version"`    // Version 状态
}

// PodStatus 定义一个结构体用于存储 Pod 的信息
type PodStatus struct {
	Name      string `json:"name"`      // Pod 名称
	Namespace string `json:"namespace"` // Pod 所属命名空间
	Status    string `json:"status"`    // Pod 状态
	NodeName  string `json:"nodename"`  // Pod 所在节点名称
}

// ListReleases 列出 Helm Releases 并存入结构体
func ListReleases(kubeconfig string) ([]Release, error) {
	// 检查 kubeconfig 文件是否存在
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Fatalf("kubeconfig 文件不存在: %v", err)
	}

	// 初始化 Helm CLI 设置
	settings := cli.New()
	settings.KubeConfig = kubeconfig

	// 初始化 Helm Action 配置
	clientGetter := kube.GetConfig(settings.KubeConfig, "", settings.Namespace()) // 指定 default 命名空间
	actionConfig := new(action.Configuration)

	err := actionConfig.Init(clientGetter, "", os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		log.Fatalf("无法初始化 Helm 配置: %v", err)
	}

	// 创建 List Client
	listClient := action.NewList(actionConfig)
	listClient.AllNamespaces = true // 查询所有命名空间

	// 查询 Releases
	releases, err := listClient.Run()
	if err != nil {
		log.Fatalf("无法获取 Releases 列表: %v", err)
	}

	// 输出结果
	fmt.Println("查询到的 Helm Releases:")
	var releaseList []Release
	if len(releases) == 0 {
		fmt.Println("没有找到已部署的 Release")
	} else {
		// 将 Release 数据存储到自定义结构体

		for _, r := range releases {

			releaseList = append(releaseList, Release{
				Name:       r.Name,
				Status:     string(r.Info.Status),
				AppVersion: r.Chart.AppVersion(),
				Version:    r.Version,
			})
		}

		// 将结构体数据转换为 JSON 格式并输出
		jsonData, err := json.MarshalIndent(releaseList, "", "  ")
		if err != nil {
			log.Fatalf("无法将 Releases 转换为 JSON: %v", err)
		}

		fmt.Println("Release 数据的 JSON 表示:")
		fmt.Println(string(jsonData))
	}
	return releaseList, nil
}

func ListPods(kubeconfig string) ([]PodStatus, error) {
	// 检查 kubeconfig 文件是否存在
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Fatalf("kubeconfig 文件不存在: %v", err)
	}

	// 加载 kubeconfig 文件
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("无法加载 kubeconfig 文件: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("无法创建 Kubernetes 客户端: %v", err)
	}

	// 获取所有命名空间的 Pods 列表
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("无法获取 Pods 列表: %v", err)
	}

	// 输出结果
	fmt.Println("查询到的 Pods 状态:")
	var podList []PodStatus
	if len(pods.Items) == 0 {
		fmt.Println("没有找到 Pod")
	} else {
		// 将 Pods 数据存储到自定义结构体
		for _, p := range pods.Items {
			podList = append(podList, PodStatus{
				Name:      p.Name,
				Namespace: p.Namespace,
				Status:    string(p.Status.Phase),
				NodeName:  p.Spec.NodeName,
			})
		}

		// 将结构体数据转换为 JSON 格式并输出
		jsonData, err := json.MarshalIndent(podList, "", "  ")
		if err != nil {
			log.Fatalf("无法将 Pods 转换为 JSON: %v", err)
		}

		fmt.Println("Pods 数据的 JSON 表示:")
		fmt.Println(string(jsonData))
	}
	return podList, nil
}
