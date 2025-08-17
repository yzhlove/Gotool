package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/package/cipher"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
	"github.com/yzhlove/Gotool/signin/server/context"
	"github.com/yzhlove/Gotool/signin/server/service/store"
)

func RegHandle(ctx *context.Context) {

	body := ctx.GetBody()
	data, err := ctx.Parse(body)
	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	reg := pb.Registry{}
	if err = proto.Unmarshal(data, &reg); err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	// 密码加强
	passwd, err := cipher.GenerateHashPasswd(reg.Passwd)
	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	// 存储加强密码
	if err = store.Save(reg.Username, passwd); err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	ctx.JSON(buildOk(nil))
}
