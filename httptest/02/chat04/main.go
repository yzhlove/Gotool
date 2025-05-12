package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log/slog"
)

// RSA 加密解密 与 OAEP 填充方案
// RSA数据的签名与验签
// RSA 加密与解密

func main() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("generate key failed", slog.Any("error", err))
		return
	}

	publicKey := &privateKey.PublicKey

	// 原始数据
	message := []byte("hello world!!!")

	// 原始数据哈希
	hashed := sha256.Sum256(message)

	// OAEP填充方案
	label := hashed[:]

	// 使用公钥对数据加密
	encryptData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, label)
	if err != nil {
		slog.Error("encrypt data failed", slog.Any("error", err))
		return
	}

	// 使用私钥对数据签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, label)
	if err != nil {
		slog.Error("encrypt signature failed", slog.Any("error", err))
		return
	}

	// 使用公钥对数据验签
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, label, signature)
	if err != nil {
		slog.Error("verify signature failed", slog.Any("error", err))
		return
	}

	// 使用私钥对数据解密
	decryptData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptData, label)
	if err != nil {
		slog.Error("decrypt data failed", slog.Any("error", err))
		return
	}

	fmt.Println(string(decryptData))

}
