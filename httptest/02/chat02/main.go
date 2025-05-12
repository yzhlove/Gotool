package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"log/slog"
)

// RSA加密解密与OAEP填充方案

func main() {

	// RSA加密解密

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("generate key error", slog.Any("error", err))
		return
	}

	publicKey := &privateKey.PublicKey
	// 消息原文
	message := []byte("hello world!")
	// 填充数据
	label := []byte("OAEP RSA Example!")

	// 使用公钥加密
	encodeStr, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, label)
	if err != nil {
		slog.Error("encrypt error", slog.Any("error", err))
		return
	}

	// 使用私钥解密
	decodeStr, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encodeStr, label)
	if err != nil {
		slog.Error("decrypt error", slog.Any("error", err))
		return
	}

	fmt.Println(string(decodeStr))
}
