package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
)

// RSA 与 ECDH

func main() {

	//aesGCM()

	ecdhTest()

}

func aesGCM() {

	key := []byte("0123456789_0123456789_0123456789")

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	plailtext := []byte("hello world")
	associatedData := []byte("associated data:example email=12234567.qq.com")

	fmt.Println("plaintext:", string(plailtext))
	fmt.Println("associatedData:", string(associatedData))

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	// 加密
	// 执行 Seal 操作：加密与认证一步到位
	// `Seal` 会将 nonce、加密后的 plaintext 和认证标签打包在一起返回。
	// 我们将 nonce 作为第一个参数传给它，它会把 nonce 附在密文的开头。
	// 也可以自己处理 nonce 的传输，第一个参数传 nil。
	ciphertext := gcm.Seal(nonce, nonce, plailtext, associatedData)

	fmt.Printf("ciphertext (HEX): %x \n\n", ciphertext)

	// 解密

	receiveNonce := ciphertext[:gcm.NonceSize()]
	encryptPayload := ciphertext[gcm.NonceSize():]

	// 执行 Open 操作：解密与认证一步到位
	// `Open` 会验证认证标签，如果认证失败，会返回一个错误。
	// 我们将 nonce 作为第一个参数传给它，它会把 nonce 附在密文的开头。
	// 也可以自己处理 nonce 的传输，第一个参数传 nil。

	plaintext, err := gcm.Open(nil, receiveNonce, encryptPayload, associatedData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("plaintext (HEX): %x \n\n", plaintext)
	fmt.Println("plaintext:", string(plaintext))

}

func ecdhTest() {

	e1 := ecdh.P256()
	priv1, err := e1.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	pub1 := priv1.PublicKey().Bytes()

	e2 := ecdh.P256()
	priv2, err := e2.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	pub2 := priv2.PublicKey().Bytes()

	// bytes 转换成pubKey
	np1, err := ecdh.P256().NewPublicKey(pub1)
	if err != nil {
		panic(err)
	}

	s1, err := priv2.ECDH(np1)
	if err != nil {
		panic(err)
	}

	np2, err := ecdh.P256().NewPublicKey(pub2)
	s2, err := priv1.ECDH(np2)
	if err != nil {
		panic(err)
	}

	fmt.Println(sha256.Sum256(s1))
	fmt.Println(sha256.Sum256(s2))

	fmt.Println("equal ==> ", hmac.Equal(s1, s2))

}
