package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func adjustKeySize(key string) []byte {
	const keySize = 32
	keyBytes := []byte(key)
	if len(keyBytes) < keySize {
		return append(keyBytes, make([]byte, keySize-len(keyBytes))...)
	}
	return keyBytes[:keySize]
}

func encrypt(data []byte) (string, error) {
	encryptSecret := os.Getenv("ENCRYPT_SECRET")
	if encryptSecret == "" {
		return "", fmt.Errorf("encrypt_secret não está definido")
	}

	block, err := aes.NewCipher(adjustKeySize(encryptSecret))
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

func decrypt(encryptedData string) ([]byte, error) {
	encryptSecret := os.Getenv("ENCRYPT_SECRET")
	if encryptSecret == "" {
		return nil, fmt.Errorf("encrypt_secret não está definido")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(adjustKeySize(encryptSecret))
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

func GenerateJWT(payload map[string]interface{}, expiryDuration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(expiryDuration)

	claims := jwt.MapClaims{
		"exp": exp.Unix(),
		"iat": now.Unix(),
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("jwt_secret não está definido")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encryptedPayload, err := encrypt(payloadBytes)
	if err != nil {
		return "", err
	}

	claims["data"] = encryptedPayload

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(adjustKeySize(secret))
	if err != nil {
		return "", err
	}

	fmt.Println(tokenString)
	return tokenString, nil
}

func DecodeJWT(tokenString string) (map[string]interface{}, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("jwt_secret não está definido")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return adjustKeySize(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if encryptedPayload, ok := claims["data"].(string); ok {
			decryptedPayload, err := decrypt(encryptedPayload)
			if err != nil {
				return nil, err
			}
			var payload map[string]interface{}
			if err := json.Unmarshal(decryptedPayload, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		}
		return nil, fmt.Errorf("dados encriptados não encontrados")
	}
	return nil, fmt.Errorf("token inválido")
}
