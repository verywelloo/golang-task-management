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
	"github.com/labstack/echo/v4"
	"github.com/patcharp/golib/cache"
	req "github.com/verywelloo/3-go-echo-task-management/app/dto/request"
	m "github.com/verywelloo/3-go-echo-task-management/app/models"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
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
	publicKeyPEM, err := AppInstance.Redis.Get(ctx, "rsa:public").Bytes()
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

func EncodeAccessToken(sessionID, userID, username string, signKey *rsa.PrivateKey) (string, error) {
	if signKey == nil {
		return "", errors.New("sign key is nil")
	}

	if userID == "" {
		return "", errors.New("subject is required")
	}

	now := time.Now()

	claim := &m.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sessionID,
			Subject:   userID, //user id
			Issuer:    "task-management",
			IssuedAt:  jwt.NewNumericDate(now),                    // NewNumericDate convert time.Time to unix timestamp
			ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour)), // toke for api-gateway service
			Audience:  []string{"api-gateway"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	return token.SignedString(signKey)
}

func LoadPublicKeyFromRedis() (*rsa.PublicKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	publicPEM, err := AppInstance.Redis.Get(ctx, "rsa:public").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("public key not found in Redis")
		}
		return nil, fmt.Errorf("redis error: %w", err)
	}

	// decode to byte
	block, _ := pem.Decode([]byte(publicPEM))
	if block == nil {
		return nil, errors.New("invalid PEM format")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("unsupported public key format")
	}

	return publicKey, nil
}

func DecodeAccessToken(accessToken string) (*m.Claims, error) {
	if accessToken == "" {
		return nil, errors.New("empty access token")
	}

	publicKey, err := LoadPublicKeyFromRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	token, err := jwt.ParseWithClaims(accessToken, &m.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*m.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
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

func GetRedis(c echo.Context, key string, result interface{}) error {
	if key == "" {
		return errors.New("redis key cannot be empty")
	}

	// get data from redis
	data, err := AppInstance.Redis.Get(c.Request().Context(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %s not found", key)
		}
		return fmt.Errorf("failed to get key %s from redis: %w", key, err)
	}

	// unmarshal
	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s %w", key, err)
	}

	return nil
}

func GetSessionCache(c echo.Context) (*req.CacheSession, error) {
	claims, err := GetAuthorizeContext(c)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("session:%s", claims.ID)

	var session req.CacheSession
	if err := Caching.Get(key, &session); err != nil {
		fmt.Printf("\nget cache error\n")
		return nil, err
	}

	return &session, nil
}

func GetAuthorizeContext(c echo.Context) (m.Claims, error) {
	//claims from the request context
	var authKey = m.ContextKey{}
	claims, ok := c.Request().Context().Value(authKey).(*m.Claims)
	if !ok || claims == nil {
		return m.Claims{}, errors.New("unauthorize claims")
	}

	return *claims, nil
}

func NewCache(cacheConfig cache.Config) cache.Redis {
	return cache.NewWithCfg(cacheConfig)
}

func VerifyPassword(candidatePassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(candidatePassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, nil
	}

	return true, nil
}
