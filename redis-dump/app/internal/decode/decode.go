package decode

import (
	"bytes"
	"encoding/json"

	"github.com/gookit/goutil"
	"github.com/vmihailenco/msgpack/v5"
)

type Type int

const (
	Unknown = iota
	JSON
	MsgPack
)

type parseFunc func([]byte) (string, bool)

var types = []Type{JSON, MsgPack}
var StringType = []string{"Unknow", "JSON", "MsgPack"}

func jsonParse(value []byte) (string, bool) {
	if json.Valid(value) {
		return string(value), true
	}
	return "", false
}

func msgPackParse(value []byte) (string, bool) {
	d := msgpack.NewDecoder(bytes.NewReader(value))
	if data, err := d.DecodeInterface(); err == nil {
		return goutil.String(data), true
	}
	return "", false
}

func Parse(value []byte) (Type, string) {
	for k, parse := range []parseFunc{jsonParse, msgPackParse} {
		if res, ok := parse(value); ok {
			return types[k], res
		}
	}
	return Unknown, ""
}
