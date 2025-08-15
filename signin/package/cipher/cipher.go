package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/yzhlove/Gotool/signin/helper"
	"golang.org/x/crypto/hkdf"
)

var errEcdsaVerify = errors.New("ecdsa verifySAN1 failed! ")

func NewDHPrivateKey() *ecdh.PrivateKey {
	return helper.Try(ecdh.P256().GenerateKey(rand.Reader)).Must()
}

func NewDHPublicKey(bytes []byte) *ecdh.PublicKey {
	return helper.Try(ecdh.P256().NewPublicKey(bytes)).Must()
}

func HKDF(secret, slot, info []byte) []byte {
	trustyKey := make([]byte, 32)
	helper.Try(io.ReadFull(hkdf.New(sha256.New, secret, slot, info), trustyKey)).Must()
	return trustyKey
}

func BuildSlot(secret string, timestamp uint64) string {
	s := strconv.FormatUint(timestamp, 10)
	var sb strings.Builder
	sb.WriteString(secret)
	sha := sha256.New()
	sha.Write([]byte(s))
	hash := hex.EncodeToString(sha.Sum(nil))
	sb.WriteString(hash)
	sb.WriteString(secret)
	return sb.String()
}

func NewAesGCM(secret []byte) cipher.AEAD {
	block := helper.Try(aes.NewCipher(secret)).Must()
	return helper.Try(cipher.NewGCM(block)).Must()
}

func GCMSeal(gcm cipher.AEAD, plaintext, additionalData []byte) []byte {
	nonce := make([]byte, gcm.NonceSize())
	helper.Try(rand.Read(nonce)).Must()
	return gcm.Seal(nonce, nonce, plaintext, additionalData)
}

func GCMOpen(gcm cipher.AEAD, ciphertext, additionalData []byte) ([]byte, error) {
	nonce := ciphertext[:gcm.NonceSize()]
	return gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], additionalData)
}

func EcdsaSignASN1(info []byte) (publicKey, signature []byte) {
	privateKey := helper.Try(ecdsa.GenerateKey(elliptic.P256(), rand.Reader)).Must()
	signature = helper.Try(ecdsa.SignASN1(rand.Reader, privateKey, info)).Must()
	pubBytes := helper.Try(privateKey.PublicKey.Bytes()).Must()
	return pubBytes, signature
}

func EcdsaVerifyASN1(publicKey, info, signature []byte) bool {
	pub := helper.Try(ecdsa.ParseUncompressedPublicKey(elliptic.P256(), publicKey)).Must()
	return ecdsa.VerifyASN1(pub, info, signature)
}

func Encode(seed uint64) (*ecdh.PrivateKey, *Meta) {
	secret := ToString(seed)
	slot := BuildSlot(secret, seed)
	info := Grow(slot, seed)

	super := HKDF([]byte(secret), []byte(slot), []byte(info))
	gcm := NewAesGCM(super)

	priv := NewDHPrivateKey()
	dhPubKey := priv.PublicKey().Bytes()

	sha := sha256.New()
	sha.Write(dhPubKey)
	sha.Write(super)
	sha.Write([]byte(info))

	spubKey, sg := EcdsaSignASN1(sha.Sum(nil))
	return priv, &Meta{
		Slot:               []byte(slot),
		Info:               []byte(info),
		DHPublicKey:        GCMSeal(gcm, dhPubKey, []byte(info)),
		SignaturePublicKey: GCMSeal(gcm, spubKey, []byte(info)),
		Signature:          GCMSeal(gcm, sg, []byte(info)),
	}
}

func Decode(seed uint64, meta *Meta) (*Meta, error) {

	secret := ToString(seed)
	slot := BuildSlot(secret, seed)
	info := Grow(slot, seed)

	super := HKDF([]byte(secret), []byte(slot), []byte(info))
	gcm := NewAesGCM(super)

	dhPubKey, err := GCMOpen(gcm, meta.DHPublicKey, []byte(info))
	if err != nil {
		return nil, err
	}

	spubKey, err := GCMOpen(gcm, meta.SignaturePublicKey, []byte(info))
	if err != nil {
		return nil, err
	}

	sg, err := GCMOpen(gcm, meta.Signature, []byte(info))
	if err != nil {
		return nil, err
	}

	sha := sha256.New()
	sha.Write(dhPubKey)
	sha.Write(super)
	sha.Write([]byte(info))

	if !EcdsaVerifyASN1(spubKey, sha.Sum(nil), sg) {
		return nil, errEcdsaVerify
	}

	return &Meta{
		Slot:               []byte(slot),
		Info:               []byte(info),
		DHPublicKey:        dhPubKey,
		Signature:          sg,
		SignaturePublicKey: spubKey,
	}, nil
}

type Meta struct {
	Slot               []byte
	Info               []byte
	DHPublicKey        []byte
	SignaturePublicKey []byte
	Signature          []byte
}
