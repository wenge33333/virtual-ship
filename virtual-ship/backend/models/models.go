package models

import (
	"database/sql"
	"time"
)

type VirtualGoodsConfig struct {
	ID                int       `db:"id" json:"id"`
	ProductID         int       `db:"product_id" json:"product_id"`
	DeliveryType      int       `db:"delivery_type" json:"delivery_type"`
	DeliveryContent   *string   `db:"delivery_content" json:"delivery_content"`
	PickRule          int       `db:"pick_rule" json:"pick_rule"`
	StockWarnThreshold int      `db:"stock_warn_threshold" json:"stock_warn_threshold"`
	IsShowStock       int       `db:"is_show_stock" json:"is_show_stock"`
	CreateTime        time.Time `db:"create_time" json:"create_time"`
	UpdateTime        time.Time `db:"update_time" json:"update_time"`
}

type CardSecret struct {
	ID            int64          `db:"id" json:"id"`
	ProductID     int            `db:"product_id" json:"product_id"`
	CardNumber    string         `db:"card_number" json:"card_number"`
	CardPassword  string         `db:"card_password" json:"card_password"`
	FaceValue     *float64       `db:"face_value" json:"face_value"`
	ExpireTime    sql.NullTime   `db:"expire_time" json:"expire_time"`
	Status        int            `db:"status" json:"status"`
	OrderID       *string        `db:"order_id" json:"order_id"`
	SellTime      sql.NullTime   `db:"sell_time" json:"sell_time"`
	ImportBatchNo *string        `db:"import_batch_no" json:"import_batch_no"`
	CreateTime    time.Time      `db:"create_time" json:"create_time"`
	UpdateTime    time.Time      `db:"update_time" json:"update_time"`
}

type CardImportLog struct {
	ID           int64     `db:"id" json:"id"`
	BatchNo      string    `db:"batch_no" json:"batch_no"`
	ProductID    int       `db:"product_id" json:"product_id"`
	FileName     string    `db:"file_name" json:"file_name"`
	TotalCount   int       `db:"total_count" json:"total_count"`
	SuccessCount int       `db:"success_count" json:"success_count"`
	FailCount    int       `db:"fail_count" json:"fail_count"`
	FailDetail   *string   `db:"fail_detail" json:"fail_detail"`
	OperatorID   int       `db:"operator_id" json:"operator_id"`
	CreateTime   time.Time `db:"create_time" json:"create_time"`
}

type DeliveryLog struct {
	ID             int64     `db:"id" json:"id"`
	OrderID        string    `db:"order_id" json:"order_id"`
	ProductID      int       `db:"product_id" json:"product_id"`
	CardID         *int64    `db:"card_id" json:"card_id"`
	DeliveryType   int       `db:"delivery_type" json:"delivery_type"`
	DeliveryContent *string  `db:"delivery_content" json:"delivery_content"`
	Status         int       `db:"status" json:"status"`
	RetryCount     int       `db:"retry_count" json:"retry_count"`
	ErrorMsg       *string   `db:"error_msg" json:"error_msg"`
	PlatformSource string    `db:"platform_source" json:"platform_source"`
	CreateTime     time.Time `db:"create_time" json:"create_time"`
	UpdateTime     time.Time `db:"update_time" json:"update_time"`
}

type PlatformConfig struct {
	ID           int        `db:"id" json:"id"`
	PlatformCode string     `db:"platform_code" json:"platform_code"`
	PlatformName string     `db:"platform_name" json:"platform_name"`
	AppKey       string     `db:"app_key" json:"app_key"`
	AppSecret    string     `db:"app_secret" json:"-"`
	SessionKey   *string    `db:"session_key" json:"-"`
	PullInterval int        `db:"pull_interval" json:"pull_interval"`
	Status       int        `db:"status" json:"status"`
	LastPullTime *time.Time `db:"last_pull_time" json:"last_pull_time"`
	CreateTime   time.Time  `db:"create_time" json:"create_time"`
	UpdateTime   time.Time  `db:"update_time" json:"update_time"`
}

type PlatformOrder struct {
	ID              int64     `db:"id" json:"id"`
	PlatformCode    string    `db:"platform_code" json:"platform_code"`
	PlatformOrderNo string    `db:"platform_order_no" json:"platform_order_no"`
	LocalOrderID    *string   `db:"local_order_id" json:"local_order_id"`
	ProductID       int       `db:"product_id" json:"product_id"`
	BuyerNick       *string   `db:"buyer_nick" json:"buyer_nick"`
	PayAmount       float64   `db:"pay_amount" json:"pay_amount"`
	PlatformStatus  string    `db:"platform_status" json:"platform_status"`
	RawData         *string   `db:"raw_data" json:"raw_data"`
	SyncTime        time.Time `db:"sync_time" json:"sync_time"`
	CreateTime      time.Time `db:"create_time" json:"create_time"`
}

type ApiCredential struct {
	ID         int64     `db:"id" json:"id"`
	AppName    string    `db:"app_name" json:"app_name"`
	ApiKey     string    `db:"api_key" json:"api_key"`
	ApiSecret  string    `db:"api_secret" json:"-"`
	RateLimit  int       `db:"rate_limit" json:"rate_limit"`
	Status     int       `db:"status" json:"status"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
	UpdateTime time.Time `db:"update_time" json:"update_time"`
}

type StoreProduct struct {
	ID          int     `db:"id" json:"id"`
	ProductType int     `db:"product_type" json:"product_type"`
}
