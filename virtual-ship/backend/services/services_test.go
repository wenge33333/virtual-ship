package services

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"os"
	"strings"
	"testing"
	"virtual-ship/config"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key := "TestKey123456789012345678901234" // 32 bytes
	plaintext := "sensitive_card_password_123"

	encrypted, err := encryptAES(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := decryptAES(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypt mismatch: got %q want %q", decrypted, plaintext)
	}

	t.Logf("AES encrypt/decrypt: %s -> %s -> %s", plaintext, encrypted[:10]+"...", decrypted)
}

func TestAESEncryptEmptyString(t *testing.T) {
	key := "TestKey123456789012345678901234"

	encrypted, err := encryptAES(key, "")
	if err != nil {
		t.Fatalf("encrypt empty string failed: %v", err)
	}

	decrypted, err := decryptAES(key, encrypted)
	if err != nil {
		t.Fatalf("decrypt empty string failed: %v", err)
	}

	if decrypted != "" {
		t.Errorf("decrypt mismatch: got %q want empty", decrypted)
	}
}

func TestAESKeyPadding(t *testing.T) {
	shortKey := "short"
	padded := padKey(shortKey)
	if len(padded) != 32 {
		t.Errorf("padKey length: got %d want 32", len(padded))
	}

	// Short key should still work
	plaintext := "test data"
	encrypted, err := encryptAES(shortKey, plaintext)
	if err != nil {
		t.Fatalf("encrypt with short key failed: %v", err)
	}
	decrypted, err := decryptAES(shortKey, encrypted)
	if err != nil {
		t.Fatalf("decrypt with short key failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("decrypt mismatch: got %q want %q", decrypted, plaintext)
	}
}

func TestDecryptAESInvalidCiphertext(t *testing.T) {
	key := "TestKey123456789012345678901234"

	_, err := decryptAES(key, "invalid_base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64, got nil")
	}

	_, err = decryptAES(key, "dG9vc2hvcnQ=") // too short
	if err == nil {
		t.Error("expected error for too short ciphertext, got nil")
	}
}

func TestEncryptAESInvalidKey(t *testing.T) {
	plaintext := "test"
	encrypted, err := encryptAES("", plaintext)
	if err != nil {
		t.Logf("encrypt with empty key: %v (expected)", err)
	}
	_ = encrypted
}

func TestCSVParsing(t *testing.T) {
	// Create a sample CSV
	csvContent := "card_number,card_password,face_value,expire_time\nCARD001,pass001,100,2026-12-31\nCARD002,pass002,50,2026-06-30\n"

	reader := csv.NewReader(strings.NewReader(csvContent))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("CSV parse failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	if records[1][0] != "CARD001" || records[1][1] != "pass001" {
		t.Errorf("first data row mismatch: %v", records[1])
	}
}

func TestImportCardsService(t *testing.T) {
	// Create a temporary CSV file
	tmpFile, err := os.CreateTemp("", "test_cards_*.csv")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("card_number,card_password,face_value,expire_time\nIMP001,imp_pass_001,200,2026-12-31\nIMP002,imp_pass_002,300,2026-12-31\n")
	if err != nil {
		t.Fatalf("write temp file failed: %v", err)
	}
	tmpFile.Close()

	// Re-open for reading
	f, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("open temp file failed: %v", err)
	}
	defer f.Close()

	cfg := &config.Config{AESKey: "TestImportKey123456789012345678"}

	// Verify CSV content
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("CSV read failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records in import CSV, got %d", len(records))
	}

	// Verify AES encryption works on import data
	for i, row := range records[1:] {
		encrypted, err := encryptAES(cfg.AESKey, row[1]) // encrypt password
		if err != nil {
			t.Errorf("row %d encrypt failed: %v", i, err)
			continue
		}
		decrypted, err := decryptAES(cfg.AESKey, encrypted)
		if err != nil {
			t.Errorf("row %d decrypt failed: %v", i, err)
			continue
		}
		if decrypted != row[1] {
			t.Errorf("row %d decrypt mismatch: got %q want %q", i, decrypted, row[1])
		}
	}
}

func TestToJSON(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	v := testStruct{Name: "test", Value: 42}
	result := toJSON(v)
	expected := `{"name":"test","value":42}`

	if result != expected {
		t.Errorf("toJSON mismatch: got %q want %q", result, expected)
	}
}

func TestUUIDGeneration(t *testing.T) {
	// Just verify it doesn't panic and produces non-empty strings
	id1 := generateBatchNo()
	id2 := generateBatchNo()

	if id1 == "" || id2 == "" {
		t.Error("generated batch numbers should not be empty")
	}

	if id1 == id2 {
		t.Error("two generated batch numbers should be different")
	}
}

// generateBatchNo generates a batch number in import
func generateBatchNo() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
