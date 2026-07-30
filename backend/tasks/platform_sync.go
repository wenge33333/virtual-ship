package tasks

import (
	"log"
	"time"
	"virtual-ship/db"
	"virtual-ship/models"
)

// PlatformSyncService 平台订单同步服务
type PlatformSyncService struct{}

// SyncPlatformOrders 同步指定平台的订单
func (s *PlatformSyncService) SyncPlatformOrders(platformCode string) error {
	if db.DB == nil {
		return nil
	}

	// 查询平台配置
	var cfg models.PlatformConfig
	err := db.DB.Get(&cfg, "SELECT * FROM eb_platform_config WHERE platform_code = ? AND status = 1", platformCode)
	if err != nil {
		return nil // 平台未启用或无配置
	}

	log.Printf("[%s] 开始拉取订单...", platformCode)

	// 根据平台代码调用对应的同步方法
	switch platformCode {
	case "taobao":
		return s.syncTaobaoOrders(&cfg)
	case "xianyu":
		return s.syncXianyuOrders(&cfg)
	case "douyin":
		return s.syncDouyinOrders(&cfg)
	case "pinduoduo":
		return s.syncPinduoduoOrders(&cfg)
	default:
		log.Printf("[%s] 不支持的平台", platformCode)
		return nil
	}
}

func (s *PlatformSyncService) syncTaobaoOrders(cfg *models.PlatformConfig) error {
	log.Printf("[taobao] 淘宝订单同步 - AppKey=%s, 拉取间隔=%ds", cfg.AppKey, cfg.PullInterval)
	s.updateLastPullTime(cfg.PlatformCode)

	// TODO: 接入淘宝开放平台 TOP SDK
	// 1. 使用 AppKey + AppSecret 获取 access_token
	// 2. 调用 taobao.trades.sold.get 拉取已售订单
	// 3. 过滤出虚拟商品订单（通过 product_type 映射）
	// 4. 写入 eb_platform_order 并触发自动发货
	return nil
}

func (s *PlatformSyncService) syncXianyuOrders(cfg *models.PlatformConfig) error {
	log.Printf("[xianyu] 闲鱼订单同步 - AppKey=%s, 拉取间隔=%ds", cfg.AppKey, cfg.PullInterval)
	s.updateLastPullTime(cfg.PlatformCode)

	// TODO: 接入闲鱼平台 API（通过淘宝开放平台的闲鱼接口）
	// 1. 使用 AppKey + AppSecret 获取 session
	// 2. 调用闲鱼订单查询接口
	// 3. 过滤待发货的虚拟商品订单
	// 4. 写入 eb_platform_order 并触发自动发货
	return nil
}

func (s *PlatformSyncService) syncDouyinOrders(cfg *models.PlatformConfig) error {
	log.Printf("[douyin] 抖音订单同步 - AppKey=%s, 拉取间隔=%ds", cfg.AppKey, cfg.PullInterval)
	s.updateLastPullTime(cfg.PlatformCode)

	// TODO: 接入抖音开放平台 API
	// 1. 使用 AppKey + AppSecret 获取 access_token
	// 2. 调用 /order/search 拉取订单列表
	// 3. 过滤虚拟商品订单
	// 4. 写入 eb_platform_order 并触发自动发货
	return nil
}

func (s *PlatformSyncService) syncPinduoduoOrders(cfg *models.PlatformConfig) error {
	log.Printf("[pinduoduo] 拼多多订单同步 - AppKey=%s, 拉取间隔=%ds", cfg.AppKey, cfg.PullInterval)
	s.updateLastPullTime(cfg.PlatformCode)

	// TODO: 接入拼多多开放平台 API
	// 1. 使用 client_id + client_secret 获取 access_token
	// 2. 调用 pdd.order.list.get 拉取订单
	// 3. 过滤虚拟商品订单
	// 4. 写入 eb_platform_order 并触发自动发货
	return nil
}

func (s *PlatformSyncService) updateLastPullTime(platformCode string) {
	db.DB.Exec("UPDATE eb_platform_config SET last_pull_time = NOW() WHERE platform_code = ?", platformCode)
}

// StartPlatformSyncTasks 启动所有平台的订单同步定时任务
func StartPlatformSyncTasks() {
	svc := &PlatformSyncService{}
	platforms := []string{"taobao", "xianyu", "douyin", "pinduoduo"}

	go func() {
		// 启动后等待 30 秒再开始首次拉取
		time.Sleep(30 * time.Second)

		for _, platform := range platforms {
			go func(code string) {
				if err := svc.SyncPlatformOrders(code); err != nil {
					log.Printf("[%s] 同步失败: %v", code, err)
				}
			}(platform)
		}

		// 之后按固定间隔拉取
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			for _, platform := range platforms {
				go func(code string) {
					if err := svc.SyncPlatformOrders(code); err != nil {
						log.Printf("[%s] 同步失败: %v", code, err)
					}
				}(platform)
			}
		}
	}()

	log.Println("平台订单同步任务已启动（每 1 分钟）")
}
