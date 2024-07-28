package jwtv2

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret      string
	EncryptSecret  string
}

type JWTGenerator struct {
	config Config
}

func New(config Config) *JWTGenerator {
	return &JWTGenerator{config: config}
}

func (j *JWTGenerator) adjustKeySize(key string) []byte {
	const keySize = 32
	keyBytes := []byte(key)
	if len(keyBytes) < keySize {
		return append(keyBytes, make([]byte, keySize-len(keyBytes))...)
	}
	return keyBytes[:keySize]
}

func (j *JWTGenerator) encrypt(data []byte) (string, error) {
	if j.config.EncryptSecret == "" {
		return "", fmt.Errorf("ENCRYPT_SECRET is not set")
	}

	block, err := aes.NewCipher(j.adjustKeySize(j.config.EncryptSecret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (j *JWTGenerator) decrypt(encryptedData string) ([]byte, error) {
	if j.config.EncryptSecret == "" {
		return nil, fmt.Errorf("ENCRYPT_SECRET is not set")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(j.adjustKeySize(j.config.EncryptSecret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (j *JWTGenerator) GenerateJWT(payload map[string]interface{}, expiryDuration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(expiryDuration)

	claims := jwt.MapClaims{
		"exp": exp.Unix(),
		"iat": now.Unix(),
	}

	if j.config.JWTSecret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encryptedPayload, err := j.encrypt(payloadBytes)
	if err != nil {
		return "", err
	}

	claims["data"] = encryptedPayload

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.adjustKeySize(j.config.JWTSecret))
	if err != nil {
		return "", err
	}

	fmt.Println(tokenString)
	return tokenString, nil
}

func (j *JWTGenerator) DecodeJWT(tokenString string) (map[string]interface{}, error) {
	if j.config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.adjustKeySize(j.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if encryptedPayload, ok := claims["data"].(string); ok {
			decryptedPayload, err := j.decrypt(encryptedPayload)
			if err != nil {
				return nil, err
			}
			var payload map[string]interface{}
			if err := json.Unmarshal(decryptedPayload, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		}
		return nil, fmt.Errorf("encrypted data not found")
	}
	return nil, fmt.Errorf("invalid token")
}
