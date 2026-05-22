package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
)

// ==================== AES 加密 ====================

type AESCrypt struct {
	key []byte
}

func NewAESCrypt(key []byte) (*AESCrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key length")
	}
	return &AESCrypt{key: key}, nil
}

func (a *AESCrypt) Encrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a *AESCrypt) Decrypt(cipherText string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("invalid ciphertext")
	}
	nonce, cipherData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, cipherData, nil)
}

// ==================== RSA 加密 ====================

type RSACrypt struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRSACryptFromPEM(privKey, pubKey []byte) (*RSACrypt, error) {
	r := &RSACrypt{}
	if privKey != nil {
		block, _ := pem.Decode(privKey)
		if block == nil {
			return nil, errors.New("invalid private key")
		}
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		}
		if err != nil {
			return nil, err
		}
		r.privateKey = privateKey.(*rsa.PrivateKey)
	}
	if pubKey != nil {
		block, _ := pem.Decode(pubKey)
		if block == nil {
			return nil, errors.New("invalid public key")
		}
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		r.publicKey = publicKey.(*rsa.PublicKey)
	}
	return r, nil
}

func GenerateRSAKeyPair(bits int) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
	pubBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}
	return pem.EncodeToMemory(privBlock), pem.EncodeToMemory(pubBlock), nil
}

func (r *RSACrypt) Encrypt(plainText []byte) (string, error) {
	if r.publicKey == nil {
		return "", errors.New("public key not set")
	}
	hash := sha256.New()
	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, r.publicKey, plainText, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (r *RSACrypt) Decrypt(cipherText string) ([]byte, error) {
	if r.privateKey == nil {
		return nil, errors.New("private key not set")
	}
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, r.privateKey, data, nil)
}

// ==================== 混合加密 ====================

func EncryptWithAES(plainText []byte, publicKey []byte) (string, string, error) {
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return "", "", err
	}
	aesCrypt, _ := NewAESCrypt(aesKey)
	encryptedText, err := aesCrypt.Encrypt(plainText)
	if err != nil {
		return "", "", err
	}
	rsaCrypt, err := NewRSACryptFromPEM(nil, publicKey)
	if err != nil {
		return "", "", err
	}
	encryptedKey, err := rsaCrypt.Encrypt(aesKey)
	if err != nil {
		return "", "", err
	}
	return encryptedText, encryptedKey, nil
}

func DecryptWithAES(encryptedText, encryptedKey string, privateKey []byte) ([]byte, error) {
	rsaCrypt, err := NewRSACryptFromPEM(privateKey, nil)
	if err != nil {
		return nil, err
	}
	aesKey, err := rsaCrypt.Decrypt(encryptedKey)
	if err != nil {
		return nil, err
	}
	aesCrypt, _ := NewAESCrypt(aesKey)
	return aesCrypt.Decrypt(encryptedText)
}
