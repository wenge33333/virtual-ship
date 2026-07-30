package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestComputeHMAC(t *testing.T) {
	secret := "my_secret"
	payload := "test_api_key1234567890"

	result := computeHMAC(secret, payload)

	// Verify manually
	expected := hmac.New(sha256.New, []byte(secret))
	expected.Write([]byte(payload))
	expectedHex := hex.EncodeToString(expected.Sum(nil))

	if result != expectedHex {
		t.Errorf("HMAC mismatch: got %q want %q", result, expectedHex)
	}
}

func TestComputeHMACDifferentPayloads(t *testing.T) {
	secret := "secret_key"

	r1 := computeHMAC(secret, "payload1")
	r2 := computeHMAC(secret, "payload2")

	if r1 == r2 {
		t.Error("HMAC for different payloads should not be equal")
	}
}

func TestComputeHMACDifferentSecrets(t *testing.T) {
	payload := "same_payload"

	r1 := computeHMAC("secret1", payload)
	r2 := computeHMAC("secret2", payload)

	if r1 == r2 {
		t.Error("HMAC with different secrets should not be equal")
	}
}

func TestAbsFunction(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tt := range tests {
		if got := abs(tt.input); got != tt.expected {
			t.Errorf("abs(%d) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestHMACAuthMissingHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		apiKey     string
		timestamp  string
		signature  string
		wantStatus int
	}{
		{"no headers", "", "", "", http.StatusOK},
		{"no timestamp", "key123", "", "", http.StatusOK},
		{"no signature", "key123", "1234567890", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", nil)
			if tt.apiKey != "" {
				c.Request.Header.Set("X-Api-Key", tt.apiKey)
			}
			if tt.timestamp != "" {
				c.Request.Header.Set("X-Timestamp", tt.timestamp)
			}
			if tt.signature != "" {
				c.Request.Header.Set("X-Signature", tt.signature)
			}

			HMACAuth()(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHMACAuthInvalidTimestamp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)
	c.Request.Header.Set("X-Api-Key", "test_key")
	c.Request.Header.Set("X-Timestamp", "invalid")
	c.Request.Header.Set("X-Signature", "some_sig")

	HMACAuth()(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}
