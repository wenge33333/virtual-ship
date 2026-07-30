package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"virtual-ship/config"
	"virtual-ship/db"
	"virtual-ship/models"
	"virtual-ship/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Cfg  *config.Config
	VG   *services.VirtualGoodsService
	Card *services.CardSecretService
	Del  *services.DeliveryService
	Plat *services.PlatformConfigService
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Code: 0, Message: "success", Data: data})
}

func fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, APIResponse{Code: code, Message: msg})
}

type ConfigReq struct {
	ProductID         int    `json:"product_id"`
	DeliveryType      int    `json:"delivery_type"`
	DeliveryContent   string `json:"delivery_content"`
	PickRule          int    `json:"pick_rule"`
	StockWarnThreshold int   `json:"stock_warn_threshold"`
	IsShowStock       int    `json:"is_show_stock"`
}

func (h *Handler) CreateConfig(c *gin.Context) {
	var req ConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if err := h.VG.CreateConfig(req.ProductID, req.DeliveryType, req.DeliveryContent, req.PickRule, req.StockWarnThreshold, req.IsShowStock); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func (h *Handler) GetConfig(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("productId"))
	cfg, err := h.VG.GetConfig(pid)
	if err != nil {
		fail(c, 404, "配置不存在")
		return
	}
	ok(c, cfg)
}

func (h *Handler) ListConfigs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.DefaultQuery("pageSize", c.DefaultQuery("limit", "20"))
	pageSize, _ := strconv.Atoi(pageSizeStr)
	result, err := h.VG.ListConfigs(page, pageSize)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *Handler) UpdateConfig(c *gin.Context) {
	var req ConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if err := h.VG.UpdateConfig(req.ProductID, req.DeliveryType, req.DeliveryContent, req.PickRule, req.StockWarnThreshold, req.IsShowStock); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func (h *Handler) DeleteConfig(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("productId"))
	if err := h.VG.DeleteConfig(pid); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func (h *Handler) ImportCards(c *gin.Context) {
	productID, _ := strconv.Atoi(c.PostForm("product_id"))
	operatorID, _ := strconv.Atoi(c.PostForm("operator_id"))

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		fail(c, 400, "请上传CSV文件")
		return
	}
	defer file.Close()

	result, err := h.Card.ImportCards(productID, operatorID, file)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *Handler) GetStockStat(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("productId"))
	stat, err := h.Card.GetStockStat(pid)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, stat)
}

func (h *Handler) GetCardList(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Query("product_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.DefaultQuery("pageSize", c.DefaultQuery("limit", "20"))
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int
	if err := db.DB.Get(&total, "SELECT COUNT(1) FROM eb_card_secret WHERE product_id = ?", pid); err != nil {
		fail(c, 500, err.Error())
		return
	}

	var cards []models.CardSecret
	if err := db.DB.Select(&cards, "SELECT * FROM eb_card_secret WHERE product_id = ? ORDER BY id DESC LIMIT ? OFFSET ?", pid, pageSize, offset); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, gin.H{"list": cards, "total": total})
}

type DeliveryReq struct {
	OrderID     string `json:"order_id"`
	ProductID   int    `json:"product_id"`
	ProductType int    `json:"product_type"`
}

func (h *Handler) AutoDelivery(c *gin.Context) {
	var req DeliveryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	result, err := h.Del.AutoDelivery(req.OrderID, req.ProductID, req.ProductType)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *Handler) ManualDelivery(c *gin.Context) {
	var req DeliveryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	result, err := h.Del.ManualDelivery(req.OrderID, req.ProductID, req.ProductType)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

type OpenAPICreateOrderReq struct {
	OrderID     string  `json:"order_id"`
	ProductID   int     `json:"product_id"`
	ProductType int     `json:"product_type"`
	PayAmount   float64 `json:"pay_amount"`
	BuyerNick   string  `json:"buyer_nick"`
	Platform    string  `json:"platform"`
	CallbackURL string  `json:"callback_url"`
	Sign        string  `json:"sign"`
	Timestamp   int64   `json:"timestamp"`
}

func (h *Handler) OpenAPICreateOrder(c *gin.Context) {
	var req OpenAPICreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}

	result, err := h.Del.AutoDelivery(req.OrderID, req.ProductID, req.ProductType)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}

	if req.CallbackURL != "" {
		go func() {
			payload, _ := json.Marshal(result)
			http.Post(req.CallbackURL, "application/json", bytes.NewReader(payload))
		}()
	}

	ok(c, result)
}

func (h *Handler) GetDeliveryLogs(c *gin.Context) {
	orderID := c.Query("order_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.DefaultQuery("pageSize", c.DefaultQuery("limit", "20"))
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	where := "WHERE 1=1"
	args := []interface{}{}
	if orderID != "" {
		where += " AND order_id = ?"
		args = append(args, orderID)
	}
	if status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}

	var total int
	countQuery := "SELECT COUNT(1) FROM eb_delivery_log " + where
	if err := db.DB.Get(&total, countQuery, args...); err != nil {
		fail(c, 500, err.Error())
		return
	}

	var logs []models.DeliveryLog
	dataQuery := "SELECT * FROM eb_delivery_log " + where + " ORDER BY create_time DESC LIMIT ? OFFSET ?"
	dataArgs := append(args, pageSize, offset)
	if err := db.DB.Select(&logs, dataQuery, dataArgs...); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, gin.H{"list": logs, "total": total})
}

func (h *Handler) RetryDelivery(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "缺少order_id参数")
		return
	}
	orderID := req.OrderID
	var log models.DeliveryLog
	if err := db.DB.Get(&log, "SELECT * FROM eb_delivery_log WHERE order_id = ? AND status = 2 AND retry_count < 3 ORDER BY create_time DESC LIMIT 1", orderID); err != nil {
		fail(c, 404, "无待重试的发货记录")
		return
	}

	result, err := h.Del.AutoDelivery(orderID, log.ProductID, log.DeliveryType)
	if err != nil {
		db.DB.Exec("UPDATE eb_delivery_log SET retry_count = retry_count + 1, error_msg = ? WHERE id = ?", err.Error(), log.ID)
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

// ========== 平台配置管理 ==========

type PlatformConfigReq struct {
	PlatformCode string `json:"platform_code" binding:"required"`
	PlatformName string `json:"platform_name" binding:"required"`
	AppKey       string `json:"app_key" binding:"required"`
	AppSecret    string `json:"app_secret" binding:"required"`
	SessionKey   string `json:"session_key"`
	PullInterval int    `json:"pull_interval"`
	Status       int    `json:"status"`
}

func (h *Handler) ListPlatformConfig(c *gin.Context) {
	configs, err := h.Plat.List()
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, configs)
}

func (h *Handler) SavePlatformConfig(c *gin.Context) {
	var req PlatformConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if req.PullInterval == 0 {
		req.PullInterval = 60
	}
	if req.Status == 0 {
		req.Status = 1
	}
	if err := h.Plat.SaveOrUpdate(req.PlatformCode, req.PlatformName, req.AppKey, req.AppSecret, req.SessionKey, req.PullInterval, req.Status); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func (h *Handler) TogglePlatformConfig(c *gin.Context) {
	var req struct {
		PlatformCode string `json:"platform_code" binding:"required"`
		Status       int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if err := h.Plat.UpdateStatus(req.PlatformCode, req.Status); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func (h *Handler) DeletePlatformConfig(c *gin.Context) {
	code := c.Param("platformCode")
	if err := h.Plat.Delete(code); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// ========== 登录桩接口（开发用，替代原 CRMEB Java 后端） ==========

func (h *Handler) AdminLogin(c *gin.Context) {
	var req struct {
		Account string `json:"account"`
		Pwd     string `json:"pwd"`
	}
	c.ShouldBindJSON(&req)

	token := "dev-token-virtual-ship-2026"
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"token": token,
			"id":    1,
		},
		"message": "success",
	})
}

func (h *Handler) AdminGetInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"account":         "admin",
			"realName":        "管理员",
			"roles":           []string{"admin"},
			"permissionsList": []string{"*:*:*"},
		},
		"message": "success",
	})
}

func (h *Handler) AdminLogout(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func (h *Handler) GetLoginPic(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"backgroundImage": "",
			"logo":            "",
			"slide":           []interface{}{},
		},
		"message": "success",
	})
}

func (h *Handler) AccountDetection(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "data": gin.H{"status": 0, "num": 0}, "message": "success"})
}

func (h *Handler) AdminGetMenus(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"data": []gin.H{
			{
				"name":      "虚拟商品",
				"component": "/virtualGoods",
				"path":      "/virtualGoods",
				"icon":      "el-icon-goods",
				"childList": []gin.H{
					{"name": "商品配置", "component": "/virtualGoods/goodsConfig", "path": "/virtualGoods/goodsConfig"},
					{"name": "卡密管理", "component": "/virtualGoods/cardManage", "path": "/virtualGoods/cardManage"},
					{"name": "发货日志", "component": "/virtualGoods/deliveryLog", "path": "/virtualGoods/deliveryLog"},
					{"name": "平台配置", "component": "/virtualGoods/platformConfig", "path": "/virtualGoods/platformConfig"},
				},
			},
		},
		"message": "success",
	})
}

// ========== 统计/版权桩接口（Dashboard 页面需要） ==========

func (h *Handler) StatisticsStub(c *gin.Context) {
	ok(c, gin.H{})
}

func (h *Handler) CopyrightInfo(c *gin.Context) {
	ok(c, gin.H{"status": 0, "copyrightContext": "", "companyImage": "", "isCopyright": false})
}
