package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"strings"
)

// AES_ECB 密码本模式，安全度极其低，永远不要使用
// "ECB企鹅"

func main() {

	aesCBC()
	fmt.Println(strings.Repeat("=", 50))
	aesCTR()
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func aesCBC() { // 密码快连接模式，加密速度慢
	// PKSC#7 填充模式
	// 密钥必须是 16，,24，,32 对应 AES128，AES192，AES256

	key := []byte("a_very_secret_key_12345678901234")
	plaintext := []byte("hello world")

	fmt.Println("plaintext:", string(plaintext))

	// 加密

	// 1. 创建一个AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize() // 16 byte for AES

	// 2. PKCS#7 填充 对明文进行填充
	paddingPlaintext := pkcs7Padding(plaintext, blockSize)

	// 3.创建密文存储区，第一个快是IV
	ciphertext := make([]byte, blockSize+len(paddingPlaintext))
	iv := ciphertext[:blockSize]

	// 4.随机生成一个IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	// 5. 创建CBC模式的加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 6. 加密
	mode.CryptBlocks(ciphertext[blockSize:], paddingPlaintext)

	fmt.Printf("ciphertext (HEX): %x \n\n", ciphertext)

	// 解密

	// 解密在另外一端
	// 1. 创建一个AES解密器
	block, err = aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// 2.从接受到密文中提取IV
	receiveIV := ciphertext[:blockSize]
	receiveCiphertext := ciphertext[blockSize:]

	// 3. 创建CBC模式的解密器
	mode = cipher.NewCBCDecrypter(block, receiveIV)

	// 4. 创建明文存储区
	decryptPaddedText := make([]byte, len(receiveCiphertext))

	// 5. 解密
	mode.CryptBlocks(decryptPaddedText, receiveCiphertext)

	// 6. 去除PKCS#7 填充
	decryptText := pkcs7UnPadding(decryptPaddedText)

	fmt.Printf("decryptText: %s \n\n", string(decryptText))
}

func aesECB() { // 电码本模式，加密速度慢 ,永远不要使用 ，安全度极其低

}

func aesCFB() { // 密文反馈模式，加密速度慢

}

func aesCTR() { // 计数器模式，加密速度快，可充分利用多核
	key := []byte("a_very_secret_key_12345678901234") // 32 byte AES 256

	plaintext := []byte("hello world abcedfg")
	fmt.Println("plaintext:", string(plaintext))

	// 加密

	// 1. 创建一个AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 创建密文存储区。与 CBC 类似，我们将 Nonce (IV) 放在密文的最前面。
	// CTR 模式的输出和输入等长，所以总长度是 Nonce 长度 + 明文长度。

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	nonce := cipherText[:aes.BlockSize]

	// 3. 随机生成一个 Nonce，保证唯一
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	// 4. 创建CTR模式的加密器
	stream := cipher.NewCTR(block, nonce)

	// 5. 加密
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	fmt.Printf("ciphertext (HEX): %x \n\n", cipherText)

	// 解密

	// 解密在另外一端
	// 1. 创建一个AES解密器
	block, err = aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 从密文中提取Nonce 和真正的密文
	receiveNonce := cipherText[:aes.BlockSize]
	receiveCipherText := cipherText[aes.BlockSize:]

	// 3. 创建CTR模式的解密器
	stream = cipher.NewCTR(block, receiveNonce)

	// 4. 创建明文的存储区域
	decryptText := make([]byte, len(receiveCipherText))

	// 5. 解密
	stream.XORKeyStream(decryptText, receiveCipherText)

	fmt.Printf("decryptText: %s \n\n", string(decryptText))
}
