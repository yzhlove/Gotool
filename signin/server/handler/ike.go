package handler

import (
	"crypto/aes"
	ccipher "crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/package/cipher"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
	"golang.org/x/crypto/hkdf"
)

func IkeHTTP(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		failAck(w, fmt.Errorf("read body error: %v", err))
		return
	}

	// 1.解析 PB 内容

	var ike = new(pb.Ike)
	if err = proto.Unmarshal(data, ike); err != nil {
		failAck(w, fmt.Errorf("unmarshal error: %v", err))
		return
	}

	// 2.检查数据是否正确
	if ike.Nonce == 0 ||
		len(ike.DHPublicKey) == 0 ||
		len(ike.EcdsaPublicKey) == 0 ||
		len(ike.Signature) == 0 {
		failAck(w, fmt.Errorf("public key is empty"))
		return
	}

	// 3.根据 NONCE 计算出 SECRET 以及 SLOT
	secret := []byte(cipher.ToString(ike.Nonce))
	slot := []byte(strconv.FormatUint(ike.Nonce, 10))

	// 4. 根据 secret 以及 slot 计算 HKDF的上下文
	info := sha256.Sum256(append(slot, secret...))

	// 5. 根据 HKDF 计算出 AES 的密钥
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdf.New(sha256.New, secret, slot, info[:]), aesKey); err != nil {
		failAck(w, fmt.Errorf("hkdf error! "))
		return
	}

	// 6. 根据 AES 解密数据
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		failAck(w, fmt.Errorf("aes new error:%v", err))
		return
	}

	gcm, err := ccipher.NewGCM(block)
	if err != nil {
		failAck(w, fmt.Errorf("aes gcm error:%v", err))
		return
	}

	// 7. 校验签名
	nonce := make([]byte, aes.BlockSize)

	gcm.Open(nil)

	dh := cipher.NewDHPrivateKey()
	secret, err := dh.ECDH(cipher.NewDHPublicKey(ike.PublicKey))
	if err != nil {
		failAck(w, fmt.Errorf("ecdh error: %v", err))
		return
	}

	seed := cipher.ToString(uint64(ike.Number))

	reader := hkdf.New(sha256.New, secret, []byte("signin"), ike.PublicKey)

	newSecert := make([]byte, 32)
	if _, err := io.ReadFull(reader, newSecert); err != nil {
		failAck(w, fmt.Errorf("hkdf error: %v", err))
		return
	}

	successAck(w, &pb.Ike{PublicKey: dh.PublicKey().Bytes()})
}
