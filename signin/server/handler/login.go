package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/helper"
	pb "github.com/yzhlove/Gotool/signin/protocol/proto"
	"github.com/yzhlove/Gotool/signin/server/context"
)

func LoginHandle(ctx *context.Context) {
	data := ctx.GetBody()
	bytes, err := ctx.Parse(data)
	if err != nil {
		ctx.JSON(buildErr(err))
		return
	}

	login := &pb.Login{}
	helper.Try(proto.Unmarshal(bytes, login)).Must()

}
