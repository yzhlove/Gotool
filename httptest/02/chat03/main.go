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

func main() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("generate private key failed", slog.Any("error", err))
		return
	}

	publicKey := &privateKey.PublicKey

	// 原文
	message := []byte("hello world!!!")

	// 使用私钥对数据签名，（一般是对数据的HASH进行签名，因为原始太长)
	hashed := sha256.Sum256(message)

	// 使用私钥对数据签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		slog.Error("rsa.SignPKCS1v15 failed", slog.Any("error", err))
		return
	}

	fmt.Println("signature => ", string(signature))

	// 使用公钥对数据验签
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		slog.Error("rsa.VerifyPKCS1v15 failed", slog.Any("error", err))
		return
	}

	fmt.Println("OK!")
}
