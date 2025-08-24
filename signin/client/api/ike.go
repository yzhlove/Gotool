package api

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/client/context"
	"github.com/yzhlove/Gotool/signin/client/http"
	"github.com/yzhlove/Gotool/signin/package/cipher"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
)

func Ike(ctx *context.Context, url string) error {

	timestamp := time.Now().Unix()
	dh, meta := cipher.Encode(uint64(timestamp))

	req := &pb.IkeReq{
		Timestamp:      uint64(timestamp),
		DHPublicKey:    meta.DHPublicKey,
		EcdsaPublicKey: meta.SignaturePublicKey,
		Signature:      meta.Signature,
	}

	bytes, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	data, err := http.Do(ctx, getIkeApi(url), &http.M{Body: bytes})
	if err != nil {
		return err
	}

	if data.Code != 0 {
		return fmt.Errorf("ike failed! reason: %s", string(data.Data))
	}

	resp := &pb.IkeResp{}
	if err := proto.Unmarshal(data.Data, resp); err != nil {
		return err
	}

	mt, err := cipher.Decode(resp.Timestamp, &cipher.Meta{
		DHPublicKey:        resp.DHPublicKey,
		SignaturePublicKey: resp.EcdsaPublicKey,
		Signature:          resp.Signature,
	})
	if err != nil {
		return err
	}

	aesKey, err := dh.ECDH(cipher.NewDHPublicKey(mt.DHPublicKey))
	if err != nil {
		return err
	}

	aesSecret := cipher.HKDF(aesKey, mt.Slot, mt.Info)

	fmt.Println("AesSecret = ", hex.EncodeToString(aesSecret))

	ctx.BindAEAD(cipher.NewAesGCM(aesSecret))
	ctx.BindToken(resp.Token)
	return nil
}
