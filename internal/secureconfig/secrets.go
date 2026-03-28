package secureconfig

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func EncryptJSON(secret []byte, value map[string]string) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return EncryptString(secret, string(bytes))
}

func DecryptJSON(secret []byte, value string) (map[string]string, error) {
	if strings.TrimSpace(value) == "" {
		return map[string]string{}, nil
	}
	plaintext, err := DecryptString(secret, value)
	if err != nil {
		return nil, err
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(plaintext), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func EncryptString(secret []byte, value string) (string, error) {
	block, err := aes.NewCipher(deriveSecretKey(secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(secret []byte, value string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(deriveSecretKey(secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", fmt.Errorf("encrypted value is too short")
	}
	nonce := raw[:gcm.NonceSize()]
	plaintext, err := gcm.Open(nil, nonce, raw[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func MaskConfig(config map[string]string) map[string]any {
	out := make(map[string]any, len(config))
	for key, value := range config {
		trimmed := strings.TrimSpace(value)
		switch {
		case trimmed == "":
			out[key] = ""
		case len(trimmed) <= 6:
			out[key] = "***"
		default:
			out[key] = trimmed[:2] + strings.Repeat("*", len(trimmed)-4) + trimmed[len(trimmed)-2:]
		}
	}
	return out
}

func deriveSecretKey(source []byte) []byte {
	sum := sha256.Sum256(source)
	return sum[:]
}
