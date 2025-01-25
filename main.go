package main

import (
	"goexample/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化 Gin 引擎
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 提供靜態資源目錄，用於提供 HTML 測試介面
	r.Static("/static", "./static")

	// 處理檔案上傳的路由
	r.PUT("/uploadConfig", func(c *gin.Context) {
		// 獲取上傳的檔案
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無法接收檔案: " + err.Error()})
			return
		}

		// 設定儲存路徑
		savePath := filepath.Join("/tmp", filepath.Base(file.Filename))

		// 將檔案儲存到指定目錄
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法儲存檔案: " + err.Error()})
			return
		}

		err = os.Setenv("KUBECONFIG", savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法設定 KUBECONFIG: " + err.Error()})
			return
		}

		// 返回成功訊息
		c.JSON(http.StatusOK, gin.H{"message": "檔案上傳成功", "path": savePath})
	})

	// 處理檔案上傳的路由
	r.PUT("/uploadRelease", func(c *gin.Context) {
		// 獲取上傳的檔案
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無法接收檔案: " + err.Error()})
			return
		}

		// 設定儲存路徑
		savePath := filepath.Join("/tmp/release", filepath.Base(file.Filename))

		result, err := utils.RunHelmDiff("test-nginx", savePath, "vscode-server", "/tmp/test.config")

		// 將檔案儲存到指定目錄
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.String(http.StatusInternalServerError, "Error")
			return
		}

		// 返回成功訊息
		c.String(http.StatusOK, result)
	})

	r.GET("/listRelease", func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		// 调用 utils.ListReleases 函数
		result, _ := utils.ListReleases("/tmp/test.config")

		// 返回 JSON 格式的响应
		c.JSON(http.StatusOK, result)
	})

	r.GET("/listPods", func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		// 调用 utils.ListReleases 函数
		result, _ := utils.ListPods("/tmp/test.config")

		// 返回 JSON 格式的响应
		c.JSON(http.StatusOK, result)
	})

	r.POST("/diffRelease", func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		// 调用 utils.ListReleases 函数
		result, err := utils.RunHelmDiff("test-nginx", "/tmp/nginx-18.3.5.tgz", "vscode-server", "/tmp/test.config")

		if err != nil {
			c.JSON(http.StatusInternalServerError, result)
			return
		}

		// 返回 JSON 格式的响应
		c.String(http.StatusOK, result)
	})

	// 啟動伺服器
	if err := r.Run(":9527"); err != nil {
		log.Fatalf("無法啟動伺服器: %v", err)
	}
}
