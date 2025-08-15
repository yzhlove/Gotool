package handler

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/package/cipher"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
	"github.com/yzhlove/Gotool/signin/server/context"
)

func IkeHandle(ctx *context.Context) {

	data := ctx.GetBody()
	ike := &pb.Ike{}
	if err := proto.Unmarshal(data, ike); err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	mt, err := cipher.Decode(ike.Timestamp, &cipher.Meta{
		DHPublicKey:        ike.DHPublicKey,
		SignaturePublicKey: ike.EcdsaPublicKey,
		Signature:          ike.Signature,
	})

	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	seed := uint64(time.Now().Unix())
	pv, pm := cipher.Encode(seed)
	aesKey, err := pv.ECDH(cipher.NewDHPublicKey(mt.DHPublicKey))
	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	newIke := &pb.Ike{
		Timestamp:      seed,
		DHPublicKey:    pm.DHPublicKey,
		EcdsaPublicKey: pm.SignaturePublicKey,
		Signature:      pm.Signature,
	}

	bytes, err := proto.Marshal(newIke)
	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	aesSecret := cipher.HKDF(aesKey, pm.Slot, pm.Info)
	// 初始化 AES-GCM
	ctx.WithAEAD(cipher.NewAesGCM(aesSecret))
	ctx.JSON(buildOk(bytes))
}
