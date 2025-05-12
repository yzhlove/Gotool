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
// RSA数据的签名与验签 PSS 模式

func main() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("generate key failed", slog.Any("error", err))
		return
	}

	publicKey := &privateKey.PublicKey

	message := []byte("hello world!!!")

	// 计算原始数据的哈希值
	hashed := sha256.Sum256(message)

	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	}

	// // 使用私钥签名
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], opts)
	if err != nil {
		slog.Error("sign failed", slog.Any("error", err))
		return
	}

	fmt.Printf("signature=> %x \n", signature)

	// 使用公钥验签
	rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, opts)
	if err != nil {
		slog.Error("verify failed", slog.Any("error", err))
		return
	}

	fmt.Println("Ok!")
}
