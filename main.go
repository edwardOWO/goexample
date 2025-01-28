package main

import (
	"archive/tar"
	"goexample/utils"
	"io"
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
	r.Static("/repo", "./static/repo")

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

		//utils.UpdateRepolist()

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
		savePath := filepath.Join("/tmp", filepath.Base(file.Filename))

		// 將檔案儲存到指定目錄
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.String(http.StatusInternalServerError, "Error while saving file")
			return
		}

		// 打開已保存的壓縮檔案
		tarFile, err := os.Open(savePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error opening tar file")
			return
		}
		defer tarFile.Close()

		// 使用 tar 解壓檔案
		tarReader := tar.NewReader(tarFile)
		extractionPath := "./static"

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break // 解壓完成
			}
			if err != nil {
				c.String(http.StatusInternalServerError, "Error reading tar file")
				return
			}

			targetPath := filepath.Join(extractionPath, header.Name)

			switch header.Typeflag {
			case tar.TypeDir: // 建立目錄
				if err := os.MkdirAll(targetPath, 0755); err != nil {
					c.String(http.StatusInternalServerError, "Error creating directory")
					return
				}
			case tar.TypeReg: // 寫入檔案
				outFile, err := os.Create(targetPath)
				if err != nil {
					c.String(http.StatusInternalServerError, "Error creating file")
					return
				}
				if _, err := io.Copy(outFile, tarReader); err != nil {
					outFile.Close()
					c.String(http.StatusInternalServerError, "Error writing file")
					return
				}
				outFile.Close()
			default:
				// 其他格式無需處理
			}
		}

		utils.UpdateRepolist("my-local-repo", "http://127.0.0.1:8888/static/repo")

		test, _ := utils.GetRepolist("my-local-repo", "http://127.0.0.1:8888/static/repo")

		// 返回成功訊息
		c.JSON(http.StatusOK, test)
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
		result, err := utils.RunHelmDiff("test1", "/tmp/release/nginx-18.3.5.tgz", "vscode-server", "/tmp/test.config")

		if err != nil {
			c.JSON(http.StatusInternalServerError, result)
			return
		}

		// 返回 JSON 格式的响应
		c.String(http.StatusOK, result)
	})

	r.POST("/installRelease", func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		result, err := utils.InstallRelease("my-local-repo", "http://127.0.0.1:8888/static/repo", "nginx", "values.yaml", "/tmp/test.config")

		if err != nil {
			c.JSON(http.StatusInternalServerError, result)
			return
		}

		// 返回 JSON 格式的响应
		c.String(http.StatusOK, result)
	})

	// 啟動伺服器
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("無法啟動伺服器: %v", err)
	}
}
