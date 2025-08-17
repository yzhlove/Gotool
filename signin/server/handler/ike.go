package handler

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/yzhlove/Gotool/signin/package/cipher"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
	"github.com/yzhlove/Gotool/signin/server/context"
	"github.com/yzhlove/Gotool/signin/server/service/manager"
)

func IkeHandle(ctx *context.Context) {

	data := ctx.GetBody()
	req := &pb.IkeReq{}
	if err := proto.Unmarshal(data, req); err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	mt, err := cipher.Decode(req.Timestamp, &cipher.Meta{
		DHPublicKey:        req.DHPublicKey,
		SignaturePublicKey: req.EcdsaPublicKey,
		Signature:          req.Signature,
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

	resp := &pb.IkeResp{
		Timestamp:      seed,
		DHPublicKey:    pm.DHPublicKey,
		EcdsaPublicKey: pm.SignaturePublicKey,
		Signature:      pm.Signature,
		Token:          uuid.New().String(),
	}

	bytes, err := proto.Marshal(resp)
	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	aesSecret := cipher.HKDF(aesKey, pm.Slot, pm.Info)
	// 初始化 AES-GCM
	ctx.Bind(cipher.NewAesGCM(aesSecret))
	// 绑定context
	manager.Bind(resp.Token, ctx)
	ctx.JSON(buildOk(bytes))
}
