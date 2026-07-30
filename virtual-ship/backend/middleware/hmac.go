package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
	"virtual-ship/db"
	"virtual-ship/models"

	"github.com/gin-gonic/gin"
)

// HMACAuth 开放 API HMAC 签名认证中间件
func HMACAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Api-Key")
		timestampStr := c.GetHeader("X-Timestamp")
		signature := c.GetHeader("X-Signature")

		if apiKey == "" || timestampStr == "" || signature == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "缺少认证参数"})
			c.Abort()
			return
		}

		// 验证时间戳有效期（5 分钟）
		ts, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "时间戳格式错误"})
			c.Abort()
			return
		}
		now := time.Now().Unix()
		if abs(now-ts) > 300 {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "请求已过期"})
			c.Abort()
			return
		}

		// 查找 API 凭据
		var cred models.ApiCredential
		if db.DB == nil {
			c.JSON(http.StatusOK, gin.H{"code": 503, "message": "服务不可用"})
			c.Abort()
			return
		}
		err = db.DB.Get(&cred, "SELECT * FROM eb_api_credential WHERE api_key = ? AND status = 1", apiKey)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "无效的 API Key"})
			c.Abort()
			return
		}

		// 验证 HMAC 签名
		payload := apiKey + timestampStr
		expected := computeHMAC(cred.ApiSecret, payload)
		if !hmac.Equal([]byte(signature), []byte(expected)) {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "签名验证失败"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func computeHMAC(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
