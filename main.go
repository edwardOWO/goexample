package main

import (
	"archive/tar"
	"fmt"
	"goexample/utils"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UpgradeRequest struct {
	ReleaseName string `json:"releasename"`
	ChartName   string `json:"chartname"`
	Namespace   string `json:"namespace"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err1 := c.Cookie("username")
		password, err2 := c.Cookie("password")

		// 如果 Cookie 不存在，或帳密不正確，則返回 401
		if err1 != nil || err2 != nil || username != "admin" || password != "password123" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort() // 阻止後續處理
			return
		}

		c.Next() // 通過驗證，繼續請求處理
	}
}
func main() {

	k8sConfig := "/tmp/config.yaml"
	repourl := "http://127.0.0.1:8888/static/repo"

	version := os.Getenv("VERSION")

	if version == "" {
		version = "test"
	}

	customer := os.Getenv("CUSTOMER")
	values := ""

	if customer != "" {
		values = fmt.Sprintf("%s.yaml", customer)
	} else {
		customer = "測試用戶"
		values = "values.yaml"
	}

	baseURL := os.Getenv("BASEURL")
	if baseURL == "" {
		baseURL = "/vscode/proxy/8888"
	}

	// 初始化 Gin 引擎
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 提供靜態資源目錄，用於提供 HTML 測試介面
	r.Static("/static", "./static")
	r.Static("/repo", "./static/repo")

	r.LoadHTMLGlob("template/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"baseURL": baseURL,
		})
	})
	r.POST("/logout", func(c *gin.Context) {
		// 移除 Cookie

		cmd := exec.Command("bash", "-c", "rm -rf /tmp/*")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("%s 失敗", string(output))
		}

		c.SetCookie("username", "", -1, "/", "", false, true)
		c.SetCookie("password", "", -1, "/", "", false, true)
		// 轉跳到 /login 頁面
		c.Redirect(http.StatusFound, "/login")
	})

	r.GET("/index", AuthMiddleware(), func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"baseURL":  baseURL,
			"customer": customer,
			"version":  version,
		})
	})

	// 登入驗證
	r.POST("/login", func(c *gin.Context) {
		// 解析請求的帳號和密碼
		var req struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入"})
			return
		}

		// 驗證帳號密碼 (這裡你可以換成查詢資料庫)
		if req.Username == "admin" && req.Password == "password123" {
			// 登入成功，設定 Cookie
			c.SetCookie("username", req.Username, 3600, "/", "", false, true)
			c.SetCookie("password", req.Password, 3600, "/", "", false, true)

			// 返回登入成功訊息
			c.JSON(http.StatusOK, gin.H{"message": "登入成功"})
		} else {
			// 登入失敗
			c.JSON(http.StatusUnauthorized, gin.H{"message": "帳號或密碼錯誤"})
		}
	})

	// 處理檔案上傳的路由
	r.PUT("/uploadConfig", AuthMiddleware(), func(c *gin.Context) {
		// 獲取上傳的檔案
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無法接收檔案: " + err.Error()})
			return
		}

		// 設定儲存路徑
		//savePath := filepath.Join("/tmp", filepath.Base(file.Filename))

		// 先儲存到 test-config.yaml
		if err := c.SaveUploadedFile(file, "/tmp/test-config.yaml"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法儲存檔案: " + err.Error()})
			return
		}

		// 檢查 test-config 狀態
		err = utils.CheckK8sConfig("/tmp/test-config.yaml")

		// 檢查異常後跳出
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "K8sConfig 檢查失敗", "path": err.Error()})
			return
		}

		// 檢查成功後除存到 /tmp/config.yaml
		if err := c.SaveUploadedFile(file, k8sConfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法儲存檔案: " + err.Error()})
			return
		}

		// 返回成功訊息
		c.JSON(http.StatusOK, gin.H{"message": "檔案上傳成功", "path": k8sConfig})
	})

	r.GET("/listRepo", AuthMiddleware(), func(c *gin.Context) {

		repolist, err := utils.GetRepolist("my-local-repo", repourl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		// 返回成功訊息
		c.JSON(http.StatusOK, repolist)

	})

	// 處理檔案上傳的路由
	r.PUT("/uploadRepo", AuthMiddleware(), func(c *gin.Context) {
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
				c.JSON(http.StatusOK, gin.H{"message": "Error reading tar file"})
				return
			}

			targetPath := filepath.Join(extractionPath, header.Name)

			switch header.Typeflag {
			case tar.TypeDir: // 建立目錄
				if err := os.MkdirAll(targetPath, 0755); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory"})
					return
				}
			case tar.TypeReg: // 寫入檔案
				outFile, err := os.Create(targetPath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating file"})
					return
				}
				if _, err := io.Copy(outFile, tarReader); err != nil {
					outFile.Close()
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing file"})
					return
				}
				outFile.Close()
			default:
				// 其他格式無需處理
			}
		}

		err = utils.UpdateRepolist("my-local-repo", repourl)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		// 返回成功訊息
		c.JSON(http.StatusOK, gin.H{"message": "repo 更新成功"})
	})

	r.GET("/listRelease", AuthMiddleware(), func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		// 调用 utils.ListReleases 函数

		configPath := k8sConfig

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
			return
		}

		result, _ := utils.ListReleases(configPath)

		// 返回 JSON 格式的响应
		c.JSON(http.StatusOK, result)
	})

	r.GET("/listPods", AuthMiddleware(), func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径
		configPath := k8sConfig

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
			return
		}
		// 调用 utils.ListReleases 函数
		result, _ := utils.ListPods(configPath)

		// 返回 JSON 格式的响应
		c.JSON(http.StatusOK, result)
	})

	r.POST("/diffRelease", AuthMiddleware(), func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		// 调用 utils.ListReleases 函数
		result, err := utils.RunHelmDiff("my-local-repo", "juiker-backend", "0.2.1", "0.2.3", k8sConfig)

		if err != nil {
			c.JSON(http.StatusInternalServerError, result)
			return
		}

		// 返回 JSON 格式的响应
		c.String(http.StatusOK, result)
	})

	r.POST("/upgradeRelease", AuthMiddleware(), func(c *gin.Context) {

		// 设置 kubeconfig 文件的路径

		var req UpgradeRequest

		// 解析 JSON 请求体
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := os.Stat(k8sConfig); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
			return
		}

		// 调用升级函数
		result, err := utils.UpgradeRelease(
			"my-local-repo",
			repourl,
			req.ReleaseName,
			req.ChartName,
			values,
			req.Namespace,
			k8sConfig,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, result)
			return
		}

		// 返回 JSON 格式的响应
		c.String(http.StatusOK, result)
	})

	r.POST("/rollbackRelease", AuthMiddleware(), func(c *gin.Context) {

		// 设置 kubeconfig 文件的路径

		var req UpgradeRequest

		// 解析 JSON 请求体
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		configPath := k8sConfig

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
			return
		}

		// 调用升级函数
		result, err := utils.RollbackRelease(
			"my-local-repo",
			repourl,
			req.ReleaseName,
			req.ChartName,
			values,
			req.Namespace,
			configPath,
		)

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
