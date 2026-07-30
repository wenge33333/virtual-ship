package config

import (
	"os"
)

type Config struct {
	ServerPort   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	RedisAddr    string
	RedisDB      int
	AESKey       string
	UploadPath   string
	JWTSecret    string
	BackupPath   string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "app"),
		DBPassword: getEnv("DB_PASSWORD", "app123"),
		DBName:     getEnv("DB_NAME", "virtual_ship"),
		RedisAddr:  getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisDB:    0,
		AESKey:     getEnv("AES_KEY", "VirtualShip2026!!"),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
		JWTSecret:  getEnv("JWT_SECRET", "virtual-ship-jwt-secret-2026"),
		BackupPath: getEnv("BACKUP_PATH", "/workspace/virtual-ship/backups"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
