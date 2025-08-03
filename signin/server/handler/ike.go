package handler

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/package/cipher"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
	"golang.org/x/crypto/hkdf"
	"io"
	"net/http"
)

func IkeHTTP(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		failAck(w, fmt.Errorf("read body error: %v", err))
		return
	}

	var ike = new(pb.Ike)
	if err = proto.Unmarshal(data, ike); err != nil {
		failAck(w, fmt.Errorf("unmarshal error: %v", err))
		return
	}

	if len(ike.PublicKey) == 0 {
		failAck(w, fmt.Errorf("public key is empty"))
		return
	}

	dh := cipher.DH()
	secret, err := dh.ECDH(cipher.NewPublicKey(ike.PublicKey))
	if err != nil {
		failAck(w, fmt.Errorf("ecdh error: %v", err))
		return
	}
	reader := hkdf.New(sha256.New, secret, []byte("signin"), ike.PublicKey)

	newSecert := make([]byte, 32)
	if _, err := io.ReadFull(reader, newSecert); err != nil {
		failAck(w, fmt.Errorf("hkdf error: %v", err))
		return
	}

	successAck(w, &pb.Ike{PublicKey: dh.PublicKey().Bytes()})
}
