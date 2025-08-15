package cipher

import (
	"crypto/ecdh"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/yzhlove/Gotool/signin/helper"
)

func Test_Echo(t *testing.T) {

	dh, req := client1()
	req = server(req)
	client2(dh, req)

}

func client1() (*ecdh.PrivateKey, *ike) {

	timestamp := time.Now().Unix()

	dh, meta := Encode(uint64(timestamp))
	ike := &ike{
		Timestamp:          uint64(timestamp),
		DHPublicKey:        meta.DHPublicKey,
		Signature:          meta.Signature,
		SignaturePublicKey: meta.SignaturePublicKey,
	}
	return dh, ike
}

func server(req *ike) *ike {

	mt, err := Decode(req.Timestamp, &Meta{
		DHPublicKey:        req.DHPublicKey,
		SignaturePublicKey: req.SignaturePublicKey,
		Signature:          req.Signature,
	})
	if err != nil {
		panic(fmt.Sprintln("parse error: ", err.Error()))
	}

	timestamp := time.Now().Unix()

	dh, mt2 := Encode(uint64(timestamp))
	aesSecret := helper.Try(dh.ECDH(NewDHPublicKey(mt.DHPublicKey))).Must()

	fmt.Println("server aes key ==> ", hex.EncodeToString(aesSecret))

	trustyAesSecret := HKDF(aesSecret, mt2.Slot, mt2.Info)
	fmt.Println("server aes2 key ==> ", hex.EncodeToString(trustyAesSecret))

	return &ike{
		Timestamp:          uint64(timestamp),
		DHPublicKey:        mt2.DHPublicKey,
		Signature:          mt2.Signature,
		SignaturePublicKey: mt2.SignaturePublicKey,
	}
}

func client2(dh *ecdh.PrivateKey, req *ike) {

	mt, err := Decode(req.Timestamp, &Meta{
		DHPublicKey:        req.DHPublicKey,
		SignaturePublicKey: req.SignaturePublicKey,
		Signature:          req.Signature,
	})

	if err != nil {
		panic(fmt.Sprintln("parse error: ", err.Error()))
	}

	aesSecret := helper.Try(dh.ECDH(NewDHPublicKey(mt.DHPublicKey))).Must()

	fmt.Println("client aes key ==> ", hex.EncodeToString(aesSecret))

	trustyAesSecret := HKDF(aesSecret, mt.Slot, mt.Info)
	fmt.Println("client aes2 key ==> ", hex.EncodeToString(trustyAesSecret))
}
