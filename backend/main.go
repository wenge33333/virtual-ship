package main

import (
	"fmt"
	"log"

	"virtual-ship/config"
	"virtual-ship/db"
	"virtual-ship/handlers"
	"virtual-ship/middleware"
	"virtual-ship/services"
	"virtual-ship/tasks"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	fmt.Println("启动虚拟商品发货系统 (Go) ...")

	if err := db.Init(cfg); err != nil {
		log.Printf("警告: 数据库连接失败: %v（将以无数据库模式运行）", err)
	}

	h := &handlers.Handler{
		Cfg:  cfg,
		VG:   &services.VirtualGoodsService{},
		Card: &services.CardSecretService{Cfg: cfg},
		Del: &services.DeliveryService{
			Cfg:  cfg,
			Card: &services.CardSecretService{Cfg: cfg},
		},
		Plat: &services.PlatformConfigService{},
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.DBHealthMiddleware())

	r.Static("/uploads", cfg.UploadPath)

	api := r.Group("/api")
	{
		// 登录接口（绕过 DB 健康检查）
		api.POST("/admin/login", h.AdminLogin)
		api.GET("/admin/getAdminInfoByToken", h.AdminGetInfo)
		api.GET("/admin/logout", h.AdminLogout)
		api.GET("/admin/getLoginPic", h.GetLoginPic)
		api.POST("/admin/login/account/detection", h.AccountDetection)
		api.GET("/admin/getMenus", h.AdminGetMenus)

		// 统计/版权桩接口
		api.GET("/admin/statistics/home/index", h.StatisticsStub)
		api.GET("/admin/statistics/home/chart/user", h.StatisticsStub)
		api.GET("/admin/statistics/home/chart/order", h.StatisticsStub)
		api.GET("/admin/statistics/home/operating/data", h.StatisticsStub)
		api.GET("/admin/copyright/get/info", h.CopyrightInfo)

		admin := api.Group("/admin/virtual")
			{
				// 商品配置
				admin.GET("/goods/config/list", h.ListConfigs)
				admin.POST("/goods/config/save", h.CreateConfig)
				admin.GET("/goods/config/:productId", h.GetConfig)
				admin.POST("/goods/config/update", h.UpdateConfig)
				admin.DELETE("/goods/config/:productId", h.DeleteConfig)
				admin.GET("/goods/config/export", h.ExportConfigs)

				// 卡密管理
				admin.POST("/card/import", h.ImportCards)
				admin.GET("/card/stock/:productId", h.GetStockStat)
				admin.GET("/card/list", h.GetCardList)
				admin.DELETE("/cards/batch", h.DeleteCards)
				admin.GET("/cards/export", h.ExportCards)

				// 发货管理
				admin.POST("/delivery/auto", h.AutoDelivery)
				admin.POST("/delivery/manual", h.ManualDelivery)
				admin.GET("/delivery/logs", h.GetDeliveryLogs)
				admin.POST("/delivery/retry", h.RetryDelivery)
				admin.DELETE("/delivery/logs/batch", h.DeleteDeliveryLogs)
				admin.GET("/delivery/logs/export", h.ExportDeliveryLogs)

				// 平台配置
				admin.GET("/platform/config/list", h.ListPlatformConfig)
				admin.POST("/platform/config/save", h.SavePlatformConfig)
				admin.POST("/platform/config/toggle", h.TogglePlatformConfig)
				admin.DELETE("/platform/config/:platformCode", h.DeletePlatformConfig)
				admin.GET("/platform/config/export", h.ExportPlatformConfigs)
			}

		// 数据备份管理
		backup := api.Group("/admin/backup")
			backup.POST("/database", h.BackupDatabase)
			backup.GET("/list", h.ListBackups)
			backup.POST("/restore", h.RestoreDatabase)
			backup.GET("/download/:filename", h.DownloadBackup)

		openAPI := api.Group("/open/v1", middleware.HMACAuth())
		{
			openAPI.POST("/delivery/create", h.OpenAPICreateOrder)
		}
	}

	addr := ":" + cfg.ServerPort
	fmt.Printf("服务器启动: http://0.0.0.0%s\n", addr)

	// Start scheduled tasks
	tasks.StartDeliveryRetryTask()
	tasks.StartStockWarnTask()
	tasks.StartPlatformSyncTasks()

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
