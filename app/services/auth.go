package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

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
	return nil
	อ่าน
	//https://chat.deepseek.com/a/chat/s/a4a8acb2-ef5c-4996-a2b0-466cb60ee6a8
	//https://chatgpt.com/c/69861c79-2328-8320-aa69-bcc8eb42b2cb
}
