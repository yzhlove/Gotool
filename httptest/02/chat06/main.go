package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
)

// 打印RSA的公钥和私钥
// 生成PEM形式的公钥和私钥

func main() {

	// 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("generate key failed", slog.Any("error", err))
		return
	}

	// 公钥
	publicKey := &privateKey.PublicKey

	// 将私钥转换为 PKCS#1 格式的字符串
	privateDerStr := x509.MarshalPKCS1PrivateKey(privateKey)
	// 将DER编码的私钥转换为PEM格式
	privatePEMStr := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateDerStr,
	})

	fmt.Println("private PEM Data: ")
	fmt.Println(string(privatePEMStr))

	// 将公钥转换为 PKCS#1 格式的字符串
	publicDerStr, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		slog.Error("marshal public key failed", slog.Any("error", err))
		return
	}

	// 将DER编码的公钥转换为PEM格式
	publicPEMStr := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicDerStr,
	})

	fmt.Println("public PEM Data: ")
	fmt.Println(string(publicPEMStr))
}
