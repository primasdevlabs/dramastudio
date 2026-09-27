package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword produces a bcrypt hash suitable for storage.
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// CheckPassword compares a plaintext password against a bcrypt hash.
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// NewAPIKey generates a prefixed API key: ds_<32 random bytes hex>.
// Only the HMAC digest is stored; the plaintext is shown once.
func NewAPIKey() (plaintext string, digest string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	plaintext = "ds_" + hex.EncodeToString(raw)
	return plaintext, APIKeyDigest(plaintext), nil
}

// APIKeyDigest is the stored representation of an API key.
func APIKeyDigest(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

// VerifyWebhookSignature validates an HMAC-SHA256 provider signature over the
// raw request body, hex or base64 encoded.
func VerifyWebhookSignature(secret string, body []byte, signature string) bool {
	if secret == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)

	sig := strings.TrimPrefix(strings.TrimSpace(signature), "sha256=")
	if decoded, err := hex.DecodeString(sig); err == nil {
		return subtle.ConstantTimeCompare(expected, decoded) == 1
	}
	if decoded, err := base64.StdEncoding.DecodeString(sig); err == nil {
		return subtle.ConstantTimeCompare(expected, decoded) == 1
	}
	return false
}

// RandomToken returns a URL-safe random token (256 bits).
func RandomToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
