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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

type UpgradeRequest struct {
	ReleaseName string `json:"releasename"`
	ChartName   string `json:"chartname"`
	Namespace   string `json:"namespace"`
}

type GetLogRequest struct {
	ReleaseName string `json:"releasename"`
	ChartName   string `json:"chartname"`
	Namespace   string `json:"namespace"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err1 := c.Cookie("username")
		password, err2 := c.Cookie("password")

		// 如果 Cookie 不存在，或帳密不正確，則返回 401
		if err1 != nil || err2 != nil || username != "admin" || password != "edward0128Juikertest123321" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort() // 阻止後續處理
			return
		}
		// 驗證 OTP

		c.Next() // 通過驗證，繼續請求處理
	}
}

type Config struct {
	K8sConfig string
	RepoURL   string
	RepoName  string
	Version   string
	Customer  string
	BaseURL   string
	Values    string
	OtpSecret string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func loadConfig() Config {
	version := getEnv("VERSION", "test")
	customer := getEnv("CUSTOMER", "測試用戶")

	values := "values"
	if customer != "測試用戶" {
		values = customer
	}

	baseURL := getEnv("BASEURL", "/vscode/proxy/8888")

	return Config{
		K8sConfig: "/tmp/config.yaml",
		RepoURL:   "http://127.0.0.1:8888/static/repo",
		RepoName:  "my-local-repo",
		Version:   version,
		Customer:  customer,
		BaseURL:   baseURL,
		Values:    values,
		OtpSecret: "JBSWY3DPEHPK3PXP",
	}
}

func otp(secret string, otpPassword string) bool {

	/*
		// Step 1: 定义 TOTP 密钥和账户信息
		issuer := "MyApp"                 // 应用名
		accountName := "user@example.com" // 账户名
		secret := "JBSWY3DPEHPK3PXP"      // Base32 编码密钥

		// Step 2: 生成 TOTP URL，遵循 otpauth 格式
		url := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, accountName, secret, issuer)

		// Step 3: 生成 QR 码并保存为图片
		err := qrcode.WriteFile(url, qrcode.Medium, 256, "qrcode.png")
		if err != nil {
			log.Fatal("生成 QR 码失败:", err)
		}

		fmt.Println("QR 码已生成并保存在 qrcode.png 文件中！")
	*/

	// **Step 2: 生成 OTP**
	otpCode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		log.Println("❌ 生成 OTP 失败:", err)
		return false
	}
	fmt.Println("📌 当前 OTP:", otpCode)

	// **Step 3: 验证 OTP**
	valid, err := totp.ValidateCustom(otpPassword, secret, time.Now(),
		totp.ValidateOpts{
			Period: 30, // 30 秒 OTP
			Skew:   1,  // 允许 1 个时间窗口偏差（±30 秒）
			Digits: 6,
		})

	if err != nil {
		log.Println("❌ OTP 验证错误:", err)
		return false
	}

	if valid {
		fmt.Println("✅ OTP 验证成功！")
		return true
	} else {
		fmt.Println("❌ OTP 验证失败！")
		return false
	}
}

func main() {

	config := loadConfig()

	// 初始化 Gin 引擎
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 提供靜態資源目錄，用於提供 HTML 測試介面
	r.Static("/static", "./static")
	r.Static("/repo", "./static/repo")
	//r.Static("/log", "/opt/log")

	r.LoadHTMLGlob("template/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"baseURL": config.BaseURL,
		})
	})
	r.POST("/logout", func(c *gin.Context) {
		// 移除 Cookie

		cmd := exec.Command("rm", "/tmp/config.yaml", "/tmp/test-config.yaml")

		log.Printf("执行命令: %s", strings.Join(cmd.Args, " "))
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("%s 失敗", string(output))
		}

		cmd = exec.Command("sh", "-c", "rm -rf /opt/log/* /opt/log/.*")
		log.Printf("執行命令: %s", strings.Join(cmd.Args, " "))
		output, err = cmd.CombinedOutput()
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
			"baseURL":  config.BaseURL,
			"customer": config.Customer,
			"version":  config.Version,
		})
	})

	// 登入驗證
	r.POST("/login", func(c *gin.Context) {
		// 解析請求的帳號和密碼
		var req struct {
			Username string `form:"username"`
			Password string `form:"password"`
			OTP      string `form:"otp"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入"})
			return
		}

		// 驗證帳號密碼 (這裡你可以換成查詢資料庫),驗證OTP
		if req.Username == "admin" && req.Password == "edward0128Juikertest123321" && otp(config.OtpSecret, req.OTP) {
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

	r.POST("/log", func(c *gin.Context) {

		var req GetLogRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 202502180000
		// 202503010000
		// 202502230000

		go utils.GetReleaseLog(req.ReleaseName, req.Namespace, "/var/log", req.StartTime, req.EndTime, config.K8sConfig)

		c.JSON(200, gin.H{"message": "開始產生 Log 紀錄"})
	})

	r.GET("/log/:filename", func(c *gin.Context) {
		// 获取路径参数
		fileName := c.Param("filename")
		filePath := "/opt/log/" + fileName

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// 读取并返回文件内容
		c.File(filePath)
	})

	r.GET("/log-check/:filename", func(c *gin.Context) {
		fileName := c.Param("filename")
		basePath := "/opt/log"
		filePath := filepath.Join(basePath, fileName+".tar.gz")
		fileDone := filepath.Join(basePath, "."+fileName+".done")
		fileRunning := filepath.Join(basePath, "."+fileName+".running")

		// 檢查 `.running` 標籤（代表任務正在執行）
		if _, err := os.Stat(fileRunning); err == nil {
			c.JSON(http.StatusAccepted, gin.H{"exists": 1}) // 202: 任務仍在執行
			return
		}

		// 檢查 `.done` 標籤（代表任務已完成）
		if _, err := os.Stat(fileDone); err == nil {
			// 檢查最終的日誌檔案是否已生成
			if _, err := os.Stat(filePath); err == nil {
				c.JSON(http.StatusOK, gin.H{"exists": 2}) // 200: 日誌檔案已生成
				return
			}
		}

		// 無相關檔案
		c.JSON(http.StatusOK, gin.H{"exists": 0}) // 200: 沒有任務記錄
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
		if err := c.SaveUploadedFile(file, config.K8sConfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法儲存檔案: " + err.Error()})
			return
		}

		// 返回成功訊息
		c.JSON(http.StatusOK, gin.H{"message": "檔案上傳成功", "path": config.BaseURL})
	})

	r.GET("/listRepo", AuthMiddleware(), func(c *gin.Context) {

		repolist, err := utils.GetRepolist(config.RepoName, config.RepoURL)

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

		err = utils.UpdateRepolist(config.RepoName, config.RepoURL)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		// 返回成功訊息
		c.JSON(http.StatusOK, gin.H{"message": "repo 更新成功"})
	})

	r.GET("/listRelease", AuthMiddleware(), func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径

		// 调用 utils.ListReleases 函数

		configPath := config.K8sConfig

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "config file not found" + err.Error()})
			return
		}

		result, _ := utils.ListReleases(configPath)

		// 返回 JSON 格式的响应
		c.JSON(http.StatusOK, result)
	})

	r.GET("/listPods", AuthMiddleware(), func(c *gin.Context) {
		// 设置 kubeconfig 文件的路径
		configPath := config.K8sConfig

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
		result, err := utils.RunHelmDiff(config.RepoName, "juiker-backend", "0.2.1", "0.2.3", config.K8sConfig)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
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

		if _, err := os.Stat(config.K8sConfig); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// 调用升级函数
		result, err := utils.UpgradeRelease(
			config.RepoName,
			config.RepoURL,
			req.ReleaseName,
			req.ChartName,
			config.Values,
			req.Namespace,
			config.K8sConfig,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
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

		configPath := config.K8sConfig

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
			return
		}

		// 调用升级函数
		result, err := utils.RollbackRelease(
			config.RepoName,
			config.RepoURL,
			req.ReleaseName,
			req.ChartName,
			config.Values,
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
