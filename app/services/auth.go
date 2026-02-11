package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	//get salt

	// hashing. default coast = 10, more is more secure, but slower
	passwordWithHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordWithHash), nil

}

func GetRSAKeys(ctx context.Context) (*rsa.PrivateKey, *rsa.PublicKey, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// get private-key PEM from redis
	privateKeyPEM, err := AppInstance.Redis.Get(ctx, "rsa:private").Bytes()
	if err == redis.Nil {
		// not found => generate new key
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, nil, err
		}

		// encode private key
		privatePEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})

		// encode public key
		publicPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		})

		// store in redis
		if err := AppInstance.Redis.Set(ctx, "rsa:private", privatePEM, 0).Err(); err != nil {
			return nil, nil, err
		}
		if err := AppInstance.Redis.Set(ctx, "rsa:public", publicPEM, 0).Err(); err != nil {
			return nil, nil, err
		}

		return privateKey, &privateKey.PublicKey, nil
	}

	if err != nil {
		return nil, nil, err
	}

	// decode private key to bytes
	privateBlock, _ := pem.Decode(privateKeyPEM)
	if privateBlock == nil {
		return nil, nil, errors.New("invalid private key PEM")
	}

	// decode private key to struct
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// get public key
	publicKeyPEM, err := AppInstance.Redis.Get(ctx, "ras:public").Bytes()
	if err != nil {
		return nil, nil, err
	}

	// decode private key to bytes
	publicBlock, _ := pem.Decode(publicKeyPEM)
	if publicBlock == nil {
		return nil, nil, errors.New("invalid public key PEM")
	}

	// decode private key to struct
	publicKey, err := x509.ParsePKCS1PublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil

}

func EncodeAccessToken(id, subject, username string, signKey *rsa.PrivateKey) (string, error) {

	if signKey == nil {
		return "", errors.New("sign key is nil")
	}

	if subject == "" {
		return "", errors.New("subject is required")
	}

	now := time.Now()

	claim := &m.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			Subject:   subject, //user id
			Issuer:    "task-management",
			IssuedAt:  jwt.NewNumericDate(now),                    // NewNumericDate convert time.Time to unix timestamp
			ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour)), // toke for api-gateway service
			Audience:  []string{"api-gateway"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	return token.SignedString(signKey)

}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}

	return hex.EncodeToString(b), nil
}

func SessionKey(sessionID string) (string, error) {
	appName := GetEnv("APP_NAME", "")
	if appName == "" {
		return "", errors.New("APP_NAME is require")
	}

	key := "session"

	return appName + ":" + key + sessionID, nil
}

func SetRedis(ctx context.Context, client *redis.Client, key string, v interface{}, timeout time.Duration) error {
	if key == "" {
		return fmt.Errorf("redis: key cannot be empty")
	}

	if timeout < 0 {
		return fmt.Errorf("redis: timeout cannot be negative")
	}

	// marshal to json
	j, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("redis: failed to marshal value to json: %w", err)
	}

	//set in redis
	cmd := client.Set(ctx, key, j, timeout)
	_, err = cmd.Result()
	if err != nil {
		return fmt.Errorf("redis: failed to set key %q: %w", key, err)
	}

	return nil
}
