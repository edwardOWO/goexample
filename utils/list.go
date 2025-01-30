package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/repo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Release 定义一个结构体用于存储 Helm Release 的信息
type Release struct {
	Name            string `json:"name"`            // Release 名称
	Status          string `json:"status"`          // Release 状态
	ChartName       string `json:"chartname"`       // Release 状态
	Namespace       string `json:"namespace"`       // Release 状态
	AppVersion      string `json:"appversion"`      // Release 状态
	Version         int    `json:"version"`         // Version 状态
	ChartVersion    string `json:"chartversion"`    // Release 状态
	NewChartVersion string `json:"newchartversion"` // Release 状态
	NeedUpdate      bool   `json:"needupdate"`
}

// PodStatus 定义一个结构体用于存储 Pod 的信息
type PodStatus struct {
	Name      string `json:"name"`      // Pod 名称
	Namespace string `json:"namespace"` // Pod 所属命名空间
	Status    string `json:"status"`    // Pod 状态
	NodeName  string `json:"nodename"`  // Pod 所在节点名称
}

type HelmRepoPackage struct {
	Name         string `json:"name"`
	ChartVersion string `json:"chartversion"`
	AppVersion   string `json:"appVersion"`
	Description  string `json:"description"`
}

// ListReleases 列出 Helm Releases 并存入结构体
func findPackageByName(packages []HelmRepoPackage, name string, currentversion string) (bool, string) {
	for _, pkg := range packages {
		if pkg.Name == name {
			if pkg.ChartVersion != currentversion {
				return true, pkg.ChartVersion
			}
		}
	}
	return false, ""
}
func ListReleases(kubeconfig string) ([]Release, error) {
	// 检查 kubeconfig 文件是否存在
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Fatalf("kubeconfig 文件不存在: %v", err)
	}

	repo, _ := GetRepolist("my-local-repo", "http://127.0.0.1:8888/static/repo")

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

			status, newversion := findPackageByName(repo, "my-local-repo"+"/"+r.Chart.Metadata.Name, r.Chart.Metadata.Version)

			releaseList = append(releaseList, Release{
				Name:            r.Name,
				Status:          string(r.Info.Status),
				ChartName:       r.Chart.Metadata.Name,
				Namespace:       r.Namespace,
				AppVersion:      r.Chart.AppVersion(),
				ChartVersion:    r.Chart.Metadata.Version,
				NewChartVersion: newversion,
				Version:         r.Version,
				NeedUpdate:      status,
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

			if string(p.Status.Phase) == "Running" {
				podList = append(podList, PodStatus{
					Name:      p.Name,
					Namespace: p.Namespace,
					Status:    string(p.Status.Phase),
					NodeName:  p.Spec.NodeName,
				})
			} else {
				podList = append(podList, PodStatus{
					Name:      p.Name,
					Namespace: p.Namespace,
					Status:    fmt.Sprintf("%s", p.Status.ContainerStatuses[0].State.Waiting.Reason),
					NodeName:  p.Spec.NodeName,
				})
			}
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

func RunHelmDiff(release, chartPath, namespace string, configPath string) (string, error) {
	// 准备 Helm diff 命令
	cmd := exec.Command("helm", "diff", "upgrade", release, chartPath, "--namespace", namespace, "--kubeconfig", configPath)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Helm diff 执行失败: %s", string(output))
		return "", err
	}

	// 打印输出
	log.Printf("Helm diff 执行成功: %s", string(output))
	return string(output), nil
}

func UpdateRepolist(repoName string, repoURL string) {
	// 設定 repository 的基本信息

	entry := &repo.Entry{
		Name: repoName,
		URL:  repoURL,
	}

	// 初始化 Helm 環境設置
	settings := cli.New()

	// 使用 getter.All 獲取所有可用的 providers
	providers := getter.All(settings)

	// 使用提供的 providers 創建 repository 實例
	repository, err := repo.NewChartRepository(entry, providers)
	if err != nil {
		log.Fatalf("Failed to create chart repository: %v", err)
	}

	// 嘗試下載 index 文件
	_, err = repository.DownloadIndexFile()
	if err != nil {
		log.Fatalf("Failed to download index file: %v", err)
	}

	// 創建 Helm repository 配置文件實例
	repoFile := repo.NewFile()

	// 添加到 repository 配置中
	repoFile.Add(entry)

	// Helm 配置文件路徑
	repoFilePath := settings.RepositoryConfig

	// 寫入配置文件
	err = repoFile.WriteFile(repoFilePath, 0644)
	if err != nil {
		log.Fatalf("Failed to write repository file: %v", err)
	}

	// 提示信息
	fmt.Printf("Successfully added repository %s with URL %s\n", repoName, repoURL)

	// 載入現有 repositories 配置文件
	repositories, err := repo.LoadFile(repoFilePath)
	if err != nil {
		log.Fatalf("Failed to load repository file: %v", err)
	}

	// 列出所有已添加的 repositories
	fmt.Println("Listing all repositories:")
	for _, repoEntry := range repositories.Repositories {
		// 使用 repo.NewChartRepository 創建 repository 實例
		repository, err := repo.NewChartRepository(repoEntry, providers)
		if err != nil {
			log.Fatalf("Failed to create chart repository for %s: %v", repoEntry.Name, err)
		}

		fmt.Printf("Name: %s, URL: %s\n", repoEntry.Name, repoEntry.URL)

		// 下載 repository 的索引文件
		_, err = repository.DownloadIndexFile()
		if err != nil {
			log.Fatalf("Failed to download index for repository %s: %v", repoEntry.Name, err)
		} else {
			fmt.Printf("Repository %s updated successfully.\n", repoEntry.Name)
		}
	}
}

func GetRepolist(repoName string, repoURL string) ([]HelmRepoPackage, error) {
	// 初始化 Helm 配置
	settings := cli.New()

	// 添加仓库到 Helm 配置
	entry := &repo.Entry{
		Name: repoName,
		URL:  repoURL,
	}
	repository, err := repo.NewChartRepository(entry, getter.All(settings))
	if err != nil {
		return nil, fmt.Errorf("failed to create chart repository: %v", err)
	}

	// 下载并更新仓库索引
	indexFile, err := repository.DownloadIndexFile()
	if err != nil {
		return nil, fmt.Errorf("failed to download index file: %v", err)
	}

	// 加载索引文件
	index, err := repo.LoadIndexFile(indexFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load index file: %v", err)
	}

	// 将包信息存入 HelmRepoPackage
	var packages []HelmRepoPackage
	for name, chartVersions := range index.Entries {
		if len(chartVersions) > 0 {
			latestVersion := chartVersions[0]
			packages = append(packages, HelmRepoPackage{
				Name:         fmt.Sprintf("%s/%s", repoName, name),
				ChartVersion: latestVersion.Version,
				AppVersion:   latestVersion.AppVersion,
				Description:  latestVersion.Description,
			})
		}
	}

	return packages, nil
}
func UpgradeRelease(repoName, repoURL, releaseName, chartName string, valuesName string, namespace string, kubeconfig string) (string, error) {

	chartPackageName, _ := GetChartPackageName(repoName, repoURL, chartName)

	// 下載檔案

	cmd := exec.Command("helm", "pull", repoName+"/"+chartName, "--kubeconfig", kubeconfig)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Helm pull 执行失败: %s", string(output))
		return "", err
	}

	// 解壓縮檔案
	cmd = exec.Command("tar", "-zxvf", chartPackageName)

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("解壓縮執行失敗: %s", string(output))
		return "", err
	}

	// 初始化 Helm 设置
	cmd = exec.Command("helm", "upgrade", releaseName, "./"+chartName, "--namespace", namespace, "-f", "./"+chartName+"/"+valuesName, "--kubeconfig", kubeconfig)

	// 执行命令
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("Helm diff 执行失败: %s", string(output))
		return "", err
	}

	// 打印输出
	log.Printf("Helm diff 执行成功: %s", string(output))

	return string(output), nil

}

func GetChartPackageName(repoName, repoURL, chartName string) (string, error) {
	// 初始化 Helm 配置
	settings := cli.New()

	// 添加仓库到 Helm 配置
	entry := &repo.Entry{
		Name: repoName,
		URL:  repoURL,
	}
	repository, err := repo.NewChartRepository(entry, getter.All(settings))
	if err != nil {
		return "", fmt.Errorf("failed to download index file: %v", err)
	}

	// 下载并更新仓库索引
	indexFile, err := repository.DownloadIndexFile()
	if err != nil {
		return "", fmt.Errorf("failed to download index file: %v", err)
	}

	// 加载索引文件
	index, err := repo.LoadIndexFile(indexFile)
	if err != nil {
		return "", fmt.Errorf("failed to load index file: %v", err)
	}

	// 检索指定 Chart 的版本信息
	chartVersions, ok := index.Entries[chartName]
	if !ok || len(chartVersions) == 0 {
		return "", fmt.Errorf("chart %s not found in repository %s", chartName, repoName)
	}

	// 返回指定 Chart 的最新版本
	latestVersion := chartVersions[0] // Helm 默认按照版本倒序排列

	packagesName := fmt.Sprintf("%s-%s.tgz", chartName, latestVersion.Version)

	return packagesName, nil
}
