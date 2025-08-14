package cipher

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strconv"
	"strings"

	"github.com/yzhlove/Gotool/signin/helper"
	"golang.org/x/crypto/hkdf"
)

func DH() *ecdh.PrivateKey {
	return helper.Try(ecdh.P256().GenerateKey(rand.Reader)).Must()
}

func NewPublicKey(bytes []byte) *ecdh.PublicKey {
	return helper.Try(ecdh.P256().NewPublicKey(bytes)).Must()
}

func HKDF(secret, slat, info []byte) []byte {
	trustySecret := make([]byte, 32)
	helper.Try(io.ReadFull(hkdf.New(sha256.New, secret, slat, info), trustySecret)).Must()
	return trustySecret
}

func BuildSlot(secret string, timestamp uint64) string {
	s := strconv.FormatUint(timestamp, 10)
	var sb strings.Builder
	sb.WriteString(hex.EncodeToString([]byte(s)))
	sb.WriteString(secret)
	return sb.String()
}
