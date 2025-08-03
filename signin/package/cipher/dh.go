package cipher

import (
	"crypto/ecdh"
	"crypto/rand"
	"github.com/yzhlove/Gotool/signin/helper"
)

func DH() *ecdh.PrivateKey {
	return helper.Try(ecdh.P256().GenerateKey(rand.Reader)).Must()
}

func NewPublicKey(bytes []byte) *ecdh.PublicKey {
	return helper.Try(ecdh.P256().NewPublicKey(bytes)).Must()
}
