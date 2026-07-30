package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"virtual-ship/db"
	"virtual-ship/models"

	"github.com/gin-gonic/gin"
)

// ========== 导出功能 ==========

func (h *Handler) ExportConfigs(c *gin.Context) {
	var configs []models.VirtualGoodsConfig
	if err := db.DB.Select(&configs, "SELECT * FROM eb_virtual_goods_config ORDER BY id"); err != nil {
		fail(c, 500, err.Error())
		return
	}

	buf := new(strings.Builder)
	buf.WriteString("\xef\xbb\xbf") // UTF-8 BOM
	w := csv.NewWriter(buf)
	w.Write([]string{"ID", "商品 ID", "发货类型", "发货内容", "选卡规则", "告警阈值", "显示库存", "创建时间"})

	for _, cfg := range configs {
		deliveryType := map[int]string{1: "自动发卡密", 2: "手动处理", 3: "固定内容"}[cfg.DeliveryType]
		pickRule := map[int]string{1: "顺序", 2: "随机"}[cfg.PickRule]
		isShow := map[int]string{1: "是", 0: "否"}[cfg.IsShowStock]
		content := ""
		if cfg.DeliveryContent != nil {
			content = *cfg.DeliveryContent
		}
		w.Write([]string{
			strconv.Itoa(cfg.ID),
			strconv.Itoa(cfg.ProductID),
			deliveryType,
			content,
			pickRule,
			strconv.Itoa(cfg.StockWarnThreshold),
			isShow,
			cfg.CreateTime.Format(time.RFC3339),
		})
	}
	w.Flush()

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=virtual_goods_config_%s.csv", time.Now().Format("20060102_150405")))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(buf.String()))
}

func (h *Handler) ExportCards(c *gin.Context) {
	pid := c.Query("product_id")
	var cards []models.CardSecret
	var err error
	if pid != "" {
		err = db.DB.Select(&cards, "SELECT * FROM eb_card_secret WHERE product_id = ? ORDER BY id", pid)
	} else {
		err = db.DB.Select(&cards, "SELECT * FROM eb_card_secret ORDER BY id")
	}
	if err != nil {
		fail(c, 500, err.Error())
		return
	}

	buf := new(strings.Builder)
	buf.WriteString("\xef\xbb\xbf") // UTF-8 BOM
	w := csv.NewWriter(buf)
	w.Write([]string{"ID", "商品 ID", "卡号", "卡密", "面值", "有效期", "状态", "订单号", "售出时间", "批次号", "创建时间"})

	for _, card := range cards {
		status := map[int]string{0: "未售", 1: "已售", 2: "已回收"}[card.Status]
		expireTime := ""
		if card.ExpireTime.Valid {
			expireTime = card.ExpireTime.Time.Format("2006-01-02")
		}
		sellTime := ""
		if card.SellTime.Valid {
			sellTime = card.SellTime.Time.Format(time.RFC3339)
		}
		orderID := ""
		if card.OrderID != nil {
			orderID = *card.OrderID
		}
		batchNo := ""
		if card.ImportBatchNo != nil {
			batchNo = *card.ImportBatchNo
		}
		w.Write([]string{
			strconv.FormatInt(card.ID, 10),
			strconv.Itoa(card.ProductID),
			card.CardNumber,
			card.CardPassword, // 已加密的密码
			fmt.Sprintf("%.2f", safeFloat(card.FaceValue)),
			expireTime,
			status,
			orderID,
			sellTime,
			batchNo,
			card.CreateTime.Format(time.RFC3339),
		})
	}
	w.Flush()

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=cards_%s.csv", time.Now().Format("20060102_150405")))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(buf.String()))
}

func (h *Handler) ExportDeliveryLogs(c *gin.Context) {
	status := c.Query("status")
	orderID := c.Query("order_id")

	where := "WHERE 1=1"
	args := []interface{}{}
	if status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}
	if orderID != "" {
		where += " AND order_id = ?"
		args = append(args, orderID)
	}

	var logs []models.DeliveryLog
	args = append(args, 1000) // LIMIT
	if err := db.DB.Select(&logs, "SELECT * FROM eb_delivery_log "+where+" ORDER BY create_time DESC LIMIT ?", args...); err != nil {
		fail(c, 500, err.Error())
		return
	}

	buf := new(strings.Builder)
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	w.Write([]string{"ID", "订单号", "商品 ID", "卡密 ID", "发货类型", "发货内容", "状态", "重试次数", "错误信息", "来源平台", "创建时间"})

	for _, log := range logs {
		deliveryType := map[int]string{1: "自动发卡密", 2: "手动处理", 3: "固定内容"}[log.DeliveryType]
		status := map[int]string{1: "成功", 2: "失败", 0: "处理中", 3: "库存不足"}[log.Status]
		content := ""
		if log.DeliveryContent != nil {
			content = *log.DeliveryContent
		}
		errorMsg := ""
		if log.ErrorMsg != nil {
			errorMsg = *log.ErrorMsg
		}
		cardID := ""
		if log.CardID != nil {
			cardID = strconv.FormatInt(*log.CardID, 10)
		}
		w.Write([]string{
			strconv.FormatInt(log.ID, 10),
			log.OrderID,
			strconv.Itoa(log.ProductID),
			cardID,
			deliveryType,
			content,
			status,
			strconv.Itoa(log.RetryCount),
			errorMsg,
			log.PlatformSource,
			log.CreateTime.Format(time.RFC3339),
		})
	}
	w.Flush()

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=delivery_logs_%s.csv", time.Now().Format("20060102_150405")))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(buf.String()))
}

func (h *Handler) ExportPlatformConfigs(c *gin.Context) {
	var configs []models.PlatformConfig
	if err := db.DB.Select(&configs, "SELECT id, platform_code, platform_name, app_key, pull_interval, status, last_pull_time, create_time, update_time FROM eb_platform_config ORDER BY id"); err != nil {
		fail(c, 500, err.Error())
		return
	}

	buf := new(strings.Builder)
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	w.Write([]string{"ID", "平台代码", "平台名称", "AppKey", "拉取间隔", "状态", "上次拉取", "创建时间"})

	for _, cfg := range configs {
		status := map[int]string{1: "启用", 0: "停用"}[cfg.Status]
		lastPull := ""
		if cfg.LastPullTime != nil {
			lastPull = cfg.LastPullTime.Format(time.RFC3339)
		}
		w.Write([]string{
			strconv.Itoa(cfg.ID),
			cfg.PlatformCode,
			cfg.PlatformName,
			cfg.AppKey,
			strconv.Itoa(cfg.PullInterval),
			status,
			lastPull,
			cfg.CreateTime.Format(time.RFC3339),
		})
	}
	w.Flush()

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=platform_configs_%s.csv", time.Now().Format("20060102_150405")))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(buf.String()))
}

// ========== 批量删除功能 ==========

func (h *Handler) DeleteCards(c *gin.Context) {
	var req struct {
		Ids []int64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if len(req.Ids) == 0 {
		fail(c, 400, "请选择要删除的卡密")
		return
	}

	placeholders := strings.Repeat("?,", len(req.Ids)-1) + "?"
	query := fmt.Sprintf("DELETE FROM eb_card_secret WHERE id IN (%s)", placeholders)
	args := make([]interface{}, len(req.Ids))
	for i, id := range req.Ids {
		args[i] = id
	}

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	n, _ := result.RowsAffected()
	ok(c, gin.H{"deleted": n})
}

func (h *Handler) DeleteDeliveryLogs(c *gin.Context) {
	var req struct {
		Ids []int64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if len(req.Ids) == 0 {
		fail(c, 400, "请选择要删除的日志")
		return
	}

	placeholders := strings.Repeat("?,", len(req.Ids)-1) + "?"
	query := fmt.Sprintf("DELETE FROM eb_delivery_log WHERE id IN (%s)", placeholders)
	args := make([]interface{}, len(req.Ids))
	for i, id := range req.Ids {
		args[i] = id
	}

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	n, _ := result.RowsAffected()
	ok(c, gin.H{"deleted": n})
}

// ========== 数据库备份恢复 ==========

func (h *Handler) BackupDatabase(c *gin.Context) {
	backupDir := h.Cfg.BackupPath
	if backupDir == "" {
		backupDir = "/workspace/virtual-ship/backups"
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fail(c, 500, "创建备份目录失败："+err.Error())
		return
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("virtual_ship_%s.sql", timestamp))

	cmd := exec.Command("mysqldump",
		"-h", "127.0.0.1",
		"-P", "3306",
		"-u", "app",
		"-papp123",
		"--databases", "virtual_ship",
		"--routines",
		"--triggers",
		"--events",
	)
	out, err := cmd.Output()
	if err != nil {
		fail(c, 500, "备份失败："+err.Error())
		return
	}

	if err := os.WriteFile(backupFile, out, 0644); err != nil {
		fail(c, 500, "保存备份文件失败："+err.Error())
		return
	}

	ok(c, gin.H{
		"filename": filepath.Base(backupFile),
		"path":     backupFile,
		"size":     len(out),
		"time":     timestamp,
	})
}

func (h *Handler) ListBackups(c *gin.Context) {
	backupDir := h.Cfg.BackupPath
	if backupDir == "" {
		backupDir = "/workspace/virtual-ship/backups"
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		fail(c, 500, "读取备份目录失败："+err.Error())
		return
	}

	var backups []gin.H
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		info, err := f.Info()
		if err != nil {
			continue
		}
		backups = append(backups, gin.H{
			"filename": f.Name(),
			"path":     filepath.Join(backupDir, f.Name()),
			"size":     info.Size(),
			"time":     info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	// 按时间倒序
	sort.Slice(backups, func(i, j int) bool {
		return backups[i]["time"].(string) > backups[j]["time"].(string)
	})

	ok(c, backups)
}

func (h *Handler) RestoreDatabase(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}

	backupDir := h.Cfg.BackupPath
	if backupDir == "" {
		backupDir = "/workspace/virtual-ship/backups"
	}

	backupFile := filepath.Join(backupDir, req.Filename)
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		fail(c, 404, "备份文件不存在")
		return
	}

	data, err := os.ReadFile(backupFile)
	if err != nil {
		fail(c, 500, "读取备份文件失败："+err.Error())
		return
	}

	cmd := exec.Command("mysql",
		"-h", "127.0.0.1",
		"-P", "3306",
		"-u", "app",
		"-papp123",
		"virtual_ship",
	)
	cmd.Stdin = strings.NewReader(string(data))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fail(c, 500, "恢复失败："+err.Error()+"\n"+string(output))
		return
	}

	ok(c, gin.H{"message": "恢复成功", "filename": req.Filename})
}

func (h *Handler) DownloadBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		fail(c, 400, "缺少文件名")
		return
	}

	backupDir := h.Cfg.BackupPath
	if backupDir == "" {
		backupDir = "/workspace/virtual-ship/backups"
	}

	backupFile := filepath.Join(backupDir, filename)
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		fail(c, 404, "文件不存在")
		return
	}

	c.Header("Content-Type", "application/sql")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(backupFile)
}

// 辅助函数
func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
