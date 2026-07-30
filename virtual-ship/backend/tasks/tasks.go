package tasks

import (
	"log"
	"time"
	"virtual-ship/db"
)

// StartDeliveryRetryTask 启动发货重试定时任务（每 5 分钟）
func StartDeliveryRetryTask() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if db.DB == nil {
				continue
			}
			// 重试失败次数小于 3 且状态为失败的发货记录
			result, err := db.DB.Exec(`
				UPDATE eb_delivery_log d
				SET d.retry_count = d.retry_count + 1, d.update_time = NOW()
				WHERE d.status = 2 AND d.retry_count < 3
				  AND d.id IN (
					SELECT id FROM (
						SELECT id FROM eb_delivery_log
						WHERE status = 2 AND retry_count < 3
						ORDER BY create_time ASC
						LIMIT 10
					) AS tmp
				  )
			`)
			if err != nil {
				log.Printf("发货重试任务执行失败: %v", err)
			} else {
				affected, _ := result.RowsAffected()
				if affected > 0 {
					log.Printf("发货重试任务: %d 条记录已标记重试", affected)
				}
			}
		}
	}()
	log.Println("发货重试定时任务已启动（每 5 分钟）")
}

// StartStockWarnTask 启动库存告警定时任务（每 5 分钟）
func StartStockWarnTask() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if db.DB == nil {
				continue
			}
			// 查询库存低于告警阈值的商品
			rows, err := db.DB.Query(`
				SELECT c.product_id, vgc.stock_warn_threshold,
				       COUNT(CASE WHEN c.status = 0 THEN 1 END) as remaining
				FROM eb_card_secret c
				JOIN eb_virtual_goods_config vgc ON c.product_id = vgc.product_id
				GROUP BY c.product_id, vgc.stock_warn_threshold
				HAVING remaining < vgc.stock_warn_threshold
			`)
			if err != nil {
				log.Printf("库存告警任务执行失败: %v", err)
				continue
			}
			count := 0
			for rows.Next() {
				var productID, threshold, remaining int
				rows.Scan(&productID, &threshold, &remaining)
				log.Printf("库存告警: 商品 %d 剩余 %d (阈值 %d)", productID, remaining, threshold)
				count++
			}
			rows.Close()
			if count > 0 {
				log.Printf("库存告警任务: 发现 %d 个商品库存低于阈值", count)
			}
		}
	}()
	log.Println("库存告警定时任务已启动（每 5 分钟）")
}
