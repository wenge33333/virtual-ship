package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"virtual-ship/config"
	"virtual-ship/db"
	"virtual-ship/models"

	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("库存不足")
	ErrConfigNotFound    = errors.New("商品虚拟配置不存在")
	ErrCardNotFound      = errors.New("卡密不存在")
	ErrDuplicateConfig   = errors.New("商品已配置虚拟发货")
)

type VirtualGoodsService struct{}

func (s *VirtualGoodsService) CreateConfig(productID, deliveryType int, content string, pickRule, warnThreshold, showStock int) error {
	var exist int
	err := db.DB.Get(&exist, "SELECT COUNT(1) FROM eb_virtual_goods_config WHERE product_id = ?", productID)
	if err != nil {
		return err
	}
	if exist > 0 {
		return ErrDuplicateConfig
	}
	_, err = db.DB.Exec(`INSERT INTO eb_virtual_goods_config 
		(product_id, delivery_type, delivery_content, pick_rule, stock_warn_threshold, is_show_stock) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		productID, deliveryType, content, pickRule, warnThreshold, showStock)
	return err
}

func (s *VirtualGoodsService) GetConfig(productID int) (*models.VirtualGoodsConfig, error) {
	var cfg models.VirtualGoodsConfig
	err := db.DB.Get(&cfg, "SELECT * FROM eb_virtual_goods_config WHERE product_id = ?", productID)
	if err != nil {
		return nil, ErrConfigNotFound
	}
	return &cfg, nil
}

func (s *VirtualGoodsService) UpdateConfig(productID, deliveryType int, content string, pickRule, warnThreshold, showStock int) error {
	_, err := db.DB.Exec(`UPDATE eb_virtual_goods_config SET 
		delivery_type=?, delivery_content=?, pick_rule=?, stock_warn_threshold=?, is_show_stock=? 
		WHERE product_id=?`,
		deliveryType, content, pickRule, warnThreshold, showStock, productID)
	return err
}

type ConfigListResult struct {
	List  []models.VirtualGoodsConfig `json:"list"`
	Total int                         `json:"total"`
}

func (s *VirtualGoodsService) ListConfigs(page, pageSize int) (*ConfigListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int
	if err := db.DB.Get(&total, "SELECT COUNT(1) FROM eb_virtual_goods_config"); err != nil {
		return nil, err
	}

	var list []models.VirtualGoodsConfig
	if err := db.DB.Select(&list, "SELECT * FROM eb_virtual_goods_config ORDER BY id DESC LIMIT ? OFFSET ?", pageSize, offset); err != nil {
		return nil, err
	}

	return &ConfigListResult{List: list, Total: total}, nil
}

func (s *VirtualGoodsService) DeleteConfig(productID int) error {
	result, err := db.DB.Exec("DELETE FROM eb_virtual_goods_config WHERE product_id = ?", productID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrConfigNotFound
	}
	return nil
}

type CardSecretService struct {
	Cfg *config.Config
}

type CardImportResult struct {
	BatchNo      string   `json:"batch_no"`
	TotalCount   int      `json:"total_count"`
	SuccessCount int      `json:"success_count"`
	FailCount    int      `json:"fail_count"`
	FailReasons  []string `json:"fail_reasons"`
}

func (s *CardSecretService) ImportCards(productID int, operatorID int, reader io.Reader) (*CardImportResult, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSV解析失败: %w", err)
	}
	if len(records) < 2 {
		return nil, errors.New("CSV文件为空或缺少数据行")
	}

	batchNo := uuid.New().String()[:16]
	result := &CardImportResult{BatchNo: batchNo, TotalCount: len(records) - 1}

	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for i, row := range records[1:] {
		if len(row) < 2 {
			result.FailCount++
			result.FailReasons = append(result.FailReasons, fmt.Sprintf("第%d行: 数据列不足", i+2))
			continue
		}
		cardNum := strings.TrimSpace(row[0])
		cardPwd := strings.TrimSpace(row[1])
		if cardNum == "" || cardPwd == "" {
			result.FailCount++
			result.FailReasons = append(result.FailReasons, fmt.Sprintf("第%d行: 卡号或密码为空", i+2))
			continue
		}

		encryptedPwd, err := encryptAES(s.Cfg.AESKey, cardPwd)
		if err != nil {
			result.FailCount++
			result.FailReasons = append(result.FailReasons, fmt.Sprintf("第%d行: 加密失败", i+2))
			continue
		}

		var faceValue *float64
		if len(row) > 2 && row[2] != "" {
			var fv float64
			fmt.Sscanf(row[2], "%f", &fv)
			faceValue = &fv
		}

		_, err = tx.Exec(`INSERT INTO eb_card_secret 
			(product_id, card_number, card_password, face_value, import_batch_no) 
			VALUES (?, ?, ?, ?, ?)`,
			productID, cardNum, encryptedPwd, faceValue, batchNo)
		if err != nil {
			result.FailCount++
			result.FailReasons = append(result.FailReasons, fmt.Sprintf("第%d行: %s(可能重复)", i+2, err.Error()))
			continue
		}
		result.SuccessCount++
	}

	_, err = tx.Exec(`INSERT INTO eb_card_import_log 
		(batch_no, product_id, file_name, total_count, success_count, fail_count, fail_detail, operator_id) 
		VALUES (?, ?, 'import.csv', ?, ?, ?, ?, ?)`,
		batchNo, productID, result.TotalCount, result.SuccessCount, result.FailCount,
		toJSON(result.FailReasons), operatorID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

type CardStockStat struct {
	Total     int `json:"total"`
	Sold      int `json:"sold"`
	Remaining int `json:"remaining"`
}

func (s *CardSecretService) GetStockStat(productID int) (*CardStockStat, error) {
	var stat CardStockStat
	err := db.DB.Get(&stat, `SELECT 
		COUNT(1) as total,
		COALESCE(SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END), 0) as sold,
		COALESCE(SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END), 0) as remaining 
		FROM eb_card_secret WHERE product_id = ?`, productID)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (s *CardSecretService) GetNextAvailable(productID int, pickRule int) (*models.CardSecret, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orderBy string
	if pickRule == 2 {
		orderBy = "RAND()"
	} else {
		orderBy = "id ASC"
	}

	var card models.CardSecret
	err = tx.Get(&card, `SELECT * FROM eb_card_secret 
		WHERE product_id = ? AND status = 0 
		ORDER BY `+orderBy+` LIMIT 1 FOR UPDATE`, productID)
	if err != nil {
		return nil, ErrInsufficientStock
	}

	_, err = tx.Exec(`UPDATE eb_card_secret SET status = 1 WHERE id = ? AND status = 0`, card.ID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &card, nil
}

func (s *CardSecretService) RecycleCard(cardID int64) error {
	_, err := db.DB.Exec(`UPDATE eb_card_secret SET status = 2 WHERE id = ? AND status = 1`, cardID)
	return err
}

type DeliveryService struct {
	Cfg  *config.Config
	Card *CardSecretService
}

type DeliveryResult struct {
	OrderID string `json:"order_id"`
	Status  int    `json:"status"`
	Content string `json:"content"`
	Msg     string `json:"msg"`
}

func (s *DeliveryService) AutoDelivery(orderID string, productID int, productType int) (*DeliveryResult, error) {
	result := &DeliveryResult{OrderID: orderID}

	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	config, err := (&VirtualGoodsService{}).GetConfig(productID)
	if err != nil {
		s.logFail(orderID, productID, 0, "", "配置查询失败: "+err.Error())
		return nil, err
	}

	switch config.DeliveryType {
	case 1:
		card, err := s.Card.GetNextAvailable(productID, config.PickRule)
		if err != nil {
			status := 0
			if err == ErrInsufficientStock {
				status = 3
			}
			s.logFail(orderID, productID, config.DeliveryType, "", err.Error())
			result.Status = status
			result.Msg = err.Error()
			return result, nil
		}

		decryptedPwd, _ := decryptAES(s.Cfg.AESKey, card.CardPassword)
		content := fmt.Sprintf("卡号: %s\n密码: %s", card.CardNumber, decryptedPwd)

		_, err = tx.Exec(`UPDATE eb_card_secret SET order_id = ?, sell_time = NOW() WHERE id = ?`,
			orderID, card.ID)
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(`INSERT INTO eb_delivery_log (order_id, product_id, card_id, delivery_type, delivery_content, status, platform_source) 
			VALUES (?, ?, ?, ?, ?, 1, 'local')`,
			orderID, productID, card.ID, config.DeliveryType, content)
		if err != nil {
			return nil, err
		}

		result.Status = 1
		result.Content = content

	case 2, 3:
		content := ""
		if config.DeliveryContent != nil {
			content = *config.DeliveryContent
		}
		_, err = tx.Exec(`INSERT INTO eb_delivery_log (order_id, product_id, delivery_type, delivery_content, status, platform_source) 
			VALUES (?, ?, ?, ?, 1, 'local')`,
			orderID, productID, config.DeliveryType, content)
		if err != nil {
			return nil, err
		}
		result.Status = 1
		result.Content = content
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	result.Msg = "发货成功"
	return result, nil
}

func (s *DeliveryService) ManualDelivery(orderID string, productID int, deliveryType int) (*DeliveryResult, error) {
	return s.AutoDelivery(orderID, productID, deliveryType)
}

func (s *DeliveryService) logFail(orderID string, productID int, deliveryType int, content string, errMsg string) {
	db.DB.Exec(`INSERT INTO eb_delivery_log (order_id, product_id, delivery_type, delivery_content, status, error_msg, platform_source) 
		VALUES (?, ?, ?, ?, 2, ?, 'local')`,
		orderID, productID, deliveryType, content, errMsg)
}

func encryptAES(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(padKey(key)))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(key, ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	block, err := aes.NewCipher([]byte(padKey(key)))
	if err != nil {
		return "", err
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return string(data), nil
}

func padKey(key string) string {
	if len(key) >= 32 {
		return key[:32]
	}
	padded := key
	for len(padded) < 32 {
		padded += "0"
	}
	return padded
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// PlatformConfigService 平台配置管理
type PlatformConfigService struct{}

func (s *PlatformConfigService) List() ([]models.PlatformConfig, error) {
	var configs []models.PlatformConfig
	err := db.DB.Select(&configs, "SELECT id, platform_code, platform_name, app_key, pull_interval, status, last_pull_time, create_time, update_time FROM eb_platform_config ORDER BY id")
	return configs, err
}

func (s *PlatformConfigService) GetByCode(code string) (*models.PlatformConfig, error) {
	var cfg models.PlatformConfig
	err := db.DB.Get(&cfg, "SELECT id, platform_code, platform_name, app_key, pull_interval, status, last_pull_time, create_time, update_time FROM eb_platform_config WHERE platform_code = ?", code)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *PlatformConfigService) SaveOrUpdate(code, name, appKey, appSecret, sessionKey string, pullInterval, status int) error {
	existing, _ := s.GetByCode(code)
	if existing != nil {
		_, err := db.DB.Exec(`UPDATE eb_platform_config SET platform_name=?, app_key=?, app_secret=?, session_key=?, pull_interval=?, status=? WHERE platform_code=?`,
			name, appKey, appSecret, sessionKey, pullInterval, status, code)
		return err
	}
	_, err := db.DB.Exec(`INSERT INTO eb_platform_config (platform_code, platform_name, app_key, app_secret, session_key, pull_interval, status) VALUES (?,?,?,?,?,?,?)`,
		code, name, appKey, appSecret, sessionKey, pullInterval, status)
	return err
}

func (s *PlatformConfigService) UpdateStatus(code string, status int) error {
	_, err := db.DB.Exec("UPDATE eb_platform_config SET status = ? WHERE platform_code = ?", status, code)
	return err
}

func (s *PlatformConfigService) Delete(code string) error {
	result, err := db.DB.Exec("DELETE FROM eb_platform_config WHERE platform_code = ?", code)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return errors.New("平台配置不存在")
	}
	return nil
}
