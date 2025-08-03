package handler

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/yzhlove/Gotool/signin/helper"
	"net/http"
)

type Msg struct {
	Code int    `json:"code,omitempty"`
	Msg  []byte `json:"msg,omitempty"`
}

func failAck(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(helper.Try(json.Marshal(&Msg{
		Code: -1,
		Msg:  []byte(err.Error()),
	})).Must())
}

func successAck(w http.ResponseWriter, msg proto.Message) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(helper.Try(json.Marshal(&Msg{
		Code: 0,
		Msg:  helper.Try(proto.Marshal(msg)).Must(),
	})).Must())
}
