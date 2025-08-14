package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand/v2"
	"strings"
	"testing"
	"time"

	"github.com/yzhlove/Gotool/signin/helper"
)

// 该函数只存在于测试用例，防止被编译进二进制包
func buildBook() {
	seed1 := make([]byte, 32)
	helper.Try(io.ReadFull(crand.Reader, seed1)).Must()

	var sb strings.Builder
	r1 := rand.New(rand.NewChaCha8(sha256.Sum256(seed1)))
	for _, idx := range r1.Perm(int(basicLen)) {
		sb.WriteRune(rune(basicText[idx]))
	}

	seedText := sb.String()
	sb.Reset()
	sb.WriteString("var seedText = \"")
	sb.WriteString(seedText)
	sb.WriteString("\"\n\n")
	sb.WriteString("var bookText = [][]byte{\n")

	clear(seed1)
	helper.Try(io.ReadFull(strings.NewReader(seedText), seed1)).Must()

	r2 := rand.New(rand.NewChaCha8(sha256.Sum256(seed1)))

	for i := 0; i < int(basicLen); i++ {
		sb.WriteString("\t[]byte(\"")
		for _, idx := range r2.Perm(int(basicLen)) {
			sb.WriteRune(rune(seedText[idx]))
		}
		sb.WriteString("\"),\n")
	}

	sb.WriteString("}\n")
	fmt.Println(sb.String())
}

func Test_GenBook(t *testing.T) {
	buildBook()
}

func Test_ToString(t *testing.T) {

	var count = 100

	for range count {
		num := rand.Uint64()
		fmt.Println("number = ", num, "\t toString = ", ToString(num))
	}

}

func Test_Grow(t *testing.T) {

	num := uint64(1234567)
	str := ToString(num)
	newStr := Grow(str, num)
	fmt.Println("number = ", num, "\t toString = ", str, "\t expend = ", newStr)

}

func Test_EncodeDecode(t *testing.T) {

	var a = uint64(1234)
	fmt.Println("number => ", a)
	var str = ToString(a)
	fmt.Println("string => ", str)
	var b, _ = ToUint64(str)
	fmt.Println("number => ", b)

	var count = 100

	for range count {
		value1 := rand.Uint64()
		str1 := ToString(value1)
		value2, _ := ToUint64(str1)
		if value1 != value2 {
			t.Error("error")
		}
	}

}

func Test_MinProcess(t *testing.T) {

	// 最小加密解密模型

	dh, ike := processByClient1()
	//fmt.Println("ike = ", ike)
	ike = processByServer(ike)

	processBYClient2(dh, ike)

}

func processByClient1() (*ecdh.PrivateKey, *ike) {

	timestamp := uint64(time.Now().Unix())

	secret := ToString(timestamp) // 原始密码本
	slat := BuildSlot(secret, timestamp)
	info := Grow(slat, timestamp)
	fmt.Println("secret = ", secret, "\t slat = ", slat, "\t info = ", info)
	newSecret := HKDF([]byte(secret), []byte(slat), []byte(info)) // 新的密钥
	fmt.Println("newSecret = ", string(newSecret))

	block := helper.Try(aes.NewCipher(newSecret)).Must()
	gcm := helper.Try(cipher.NewGCM(block)).Must() // AES-GCM

	dh := helper.Try(ecdh.P256().GenerateKey(crand.Reader)).Must()
	publicKey := dh.PublicKey().Bytes()

	nonce := make([]byte, gcm.NonceSize())
	helper.Try(crand.Read(nonce)).Must()

	// 需要发送的公钥
	encryptPublicKey := gcm.Seal(nonce, nonce, publicKey, []byte(info)) // 加密之后的公钥

	// 签名
	privKey := helper.Try(ecdsa.GenerateKey(elliptic.P256(), crand.Reader)).Must() // ECDSA

	sg := helper.Try(ecdsa.SignASN1(crand.Reader, privKey, append(publicKey, newSecret...))).Must()
	sgPubBytes := helper.Try(privKey.PublicKey.Bytes()).Must()
	encryptSgPublicKey := gcm.Seal(nonce, nonce, sgPubBytes, []byte(info)) // 加密之后的签名
	encryptSG := gcm.Seal(nonce, nonce, sg, []byte(info))                  // 加密之后的签名

	return dh, &ike{
		Timestamp:          timestamp,
		DHPublicKey:        encryptPublicKey,
		Signature:          encryptSG,
		SignaturePublicKey: encryptSgPublicKey,
	}
}

func processBYClient2(dh *ecdh.PrivateKey, req *ike) {

	secret := ToString(req.Timestamp)
	slot := BuildSlot(secret, req.Timestamp)
	info := Grow(slot, req.Timestamp)

	newSecret := HKDF([]byte(secret), []byte(slot), []byte(info))

	block := helper.Try(aes.NewCipher(newSecret)).Must()
	gcm := helper.Try(cipher.NewGCM(block)).Must() // AES-GCM

	// 解签
	decryptSg := helper.Try(gcm.Open(nil, req.Signature[:gcm.NonceSize()], req.Signature[gcm.NonceSize():], []byte(info))).Must()
	decryptSgPub := helper.Try(gcm.Open(nil, req.SignaturePublicKey[:gcm.NonceSize()], req.SignaturePublicKey[gcm.NonceSize():], []byte(info))).Must()
	pubKey := helper.Try(ecdsa.ParseUncompressedPublicKey(elliptic.P256(), decryptSgPub)).Must()

	// 解密 DH
	decryptPublicKey := helper.Try(gcm.Open(nil, req.DHPublicKey[:gcm.NonceSize()], req.DHPublicKey[gcm.NonceSize():], []byte(info))).Must()

	// 验签
	if !ecdsa.VerifyASN1(pubKey, append(decryptPublicKey, newSecret...), decryptSg) {
		panic("verify error")
	} else {
		fmt.Println("verify Ok!!! ")
	}

	// 生成 DH
	parse := helper.Try(ecdh.P256().NewPublicKey(decryptPublicKey)).Must()
	aesSecret := helper.Try(dh.ECDH(parse)).Must()

	fmt.Println("new dh key ==> ", hex.EncodeToString(aesSecret))

	newAesKey := HKDF(aesSecret, []byte(slot), []byte(info))
	fmt.Println("new aes key ==> ", hex.EncodeToString(newAesKey))

}

func processByServer(req *ike) *ike {

	secret := ToString(req.Timestamp)
	slot := BuildSlot(secret, req.Timestamp)
	info := Grow(slot, req.Timestamp)

	newSecret := HKDF([]byte(secret), []byte(slot), []byte(info))

	block := helper.Try(aes.NewCipher(newSecret)).Must()
	gcm := helper.Try(cipher.NewGCM(block)).Must() // AES-GCM

	// 解签
	decryptSg := helper.Try(gcm.Open(nil, req.Signature[:gcm.NonceSize()], req.Signature[gcm.NonceSize():], []byte(info))).Must()
	decryptSgPub := helper.Try(gcm.Open(nil, req.SignaturePublicKey[:gcm.NonceSize()], req.SignaturePublicKey[gcm.NonceSize():], []byte(info))).Must()
	pubKey := helper.Try(ecdsa.ParseUncompressedPublicKey(elliptic.P256(), decryptSgPub)).Must()

	// 解密 DH
	decryptPublicKey := helper.Try(gcm.Open(nil, req.DHPublicKey[:gcm.NonceSize()], req.DHPublicKey[gcm.NonceSize():], []byte(info))).Must()

	// 验签
	if !ecdsa.VerifyASN1(pubKey, append(decryptPublicKey, newSecret...), decryptSg) {
		panic("verify error")
	} else {
		fmt.Println("verify Ok!!! ")
	}

	// 生成 DH
	dh := helper.Try(ecdh.P256().GenerateKey(crand.Reader)).Must()
	parse := helper.Try(ecdh.P256().NewPublicKey(decryptPublicKey)).Must()
	aesSecret := helper.Try(dh.ECDH(parse)).Must()

	fmt.Println("new dh key ==> ", hex.EncodeToString(aesSecret))

	// 新的加密流程
	timestamp := uint64(time.Now().Unix())
	secret2 := ToString(timestamp) // 原始密码本
	slot2 := BuildSlot(secret, timestamp)
	info2 := Grow(slot2, timestamp)
	newSecret2 := HKDF([]byte(secret2), []byte(slot2), []byte(info2)) // 新的密钥

	newAesKey := HKDF(aesSecret, []byte(slot2), []byte(info2))
	fmt.Println("new aes key ==> ", hex.EncodeToString(newAesKey))

	{

		block2 := helper.Try(aes.NewCipher(newSecret2)).Must()
		gcm2 := helper.Try(cipher.NewGCM(block2)).Must() // AES-GCM

		nonce := make([]byte, gcm.NonceSize())
		helper.Try(crand.Read(nonce)).Must()

		ppub := dh.PublicKey().Bytes()

		// 需要发送的公钥
		encryptPublicKey2 := gcm2.Seal(nonce, nonce, ppub, []byte(info2)) // 加密之后的公钥

		// 签名
		privKey2 := helper.Try(ecdsa.GenerateKey(elliptic.P256(), crand.Reader)).Must() // ECDSA

		sg2 := helper.Try(ecdsa.SignASN1(crand.Reader, privKey2, append(ppub, newSecret2...))).Must()
		sgPubBytes2 := helper.Try(privKey2.PublicKey.Bytes()).Must()
		encryptSgPublicKey2 := gcm.Seal(nonce, nonce, sgPubBytes2, []byte(info)) // 加密之后的签名
		encryptSG2 := gcm.Seal(nonce, nonce, sg2, []byte(info))                  // 加密之后的签名

		return &ike{
			Timestamp:          timestamp,
			DHPublicKey:        encryptPublicKey2,
			Signature:          encryptSG2,
			SignaturePublicKey: encryptSgPublicKey2,
		}
	}

	return nil

}

type ike struct {
	Timestamp          uint64
	DHPublicKey        []byte
	Signature          []byte
	SignaturePublicKey []byte
}
