package db

import (
	"fmt"
	"virtual-ship/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(10)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
