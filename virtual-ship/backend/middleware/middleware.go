package middleware

import (
	"net/http"
	"virtual-ship/db"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Api-Key, X-Timestamp, X-Signature")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func DBHealthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if db.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "数据库未连接，服务不可用"})
			c.Abort()
			return
		}
		if err := db.DB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "数据库连接异常"})
			c.Abort()
			return
		}
		c.Next()
	}
}
