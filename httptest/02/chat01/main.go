package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log/slog"
)

// ECDH使用

func main() {

	// 私密key

	alicePrivKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		slog.Error("ecdh.P256().GenerateKey failed", slog.Any("error", err))
		return
	}

	bobPrivKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		slog.Error("ecdh.P256().GenerateKey failed", slog.Any("error", err))
		return
	}

	fmt.Printf("alicePrivKey => %x\n", sha256.Sum256(alicePrivKey.Bytes()))
	fmt.Printf("bobPrivKey => %x\n", sha256.Sum256(bobPrivKey.Bytes()))

	// 公钥
	alicePubKey := alicePrivKey.PublicKey()
	bobPubKey := bobPrivKey.PublicKey()

	// alice和bob交换公钥，然后使用公钥进行ECDH交换
	secret1, err := alicePrivKey.ECDH(bobPubKey)
	if err != nil {
		slog.Error("alicePrivKey.ECDH failed", slog.Any("error", err))
		return
	}

	// bob和alice交换公钥，然后使用公钥进行ECDH交换
	secret2, err := bobPrivKey.ECDH(alicePubKey)
	if err != nil {
		slog.Error("bobPrivKey.ECDH failed", slog.Any("error", err))
		return
	}

	fmt.Printf("secret1 => %x\n", sha256.Sum256(secret1))
	fmt.Printf("secret2 => %x\n", sha256.Sum256(secret2))

	// 验证secret1和secret2是否相等
	if !bytes.Equal(secret1, secret2) {
		slog.Error("secret1 and secret2 are not equal")
		return
	} else {
		fmt.Println("Ok")
	}
}
