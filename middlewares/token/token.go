package middlewares

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(ttl time.Duration, payload interface{}, privateKeyPath string) (string, error) {
	// Load private key from file
	key, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("could not load private key: %w", err)
	}

	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub": payload,             // User ID atau data lainnya
		"exp": now.Add(ttl).Unix(), // Expiry time
		"iat": now.Unix(),          // Issued at
		"nbf": now.Unix(),          // Not before
	}

	// Buat token dengan RSA256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return signedToken, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(path) // ✅ Gunakan os.ReadFile
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}
	return privateKey, nil
}

// Load public key from file
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(path) // ✅ Gunakan os.ReadFile
	if err != nil {
		return nil, fmt.Errorf("error reading public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}
	return publicKey, nil
}
