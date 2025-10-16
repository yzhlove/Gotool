package rdb

import (
	"strings"
	"time"

	"github.com/hdt3213/rdb/model"
	"rain.com/Gotool/redis-dump/app/internal/decode"
)

type Meta struct {
	Values    ListData
	RedisData model.RedisObject
}

type Data struct {
	Value string
	Name  string
	Type  decode.Type
}

type ListData []Data

func (ls ListData) String() string {
	var sb strings.Builder
	sb.Grow(64)
	for _, v := range ls {
		if len(v.Name) > 0 {
			sb.WriteString("Name:")
			sb.WriteString(v.Name)
			sb.WriteString(",")
		}
		sb.WriteString("Type:")
		sb.WriteString(decode.StringType[v.Type])
		sb.WriteString(",")
		sb.WriteString("Value:")
		sb.WriteString(v.Value)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m Meta) String() string {
	var sb strings.Builder
	sb.Grow(64)

	sb.WriteString("Key: ")
	sb.WriteString(m.RedisData.GetKey())
	sb.WriteString(" ")
	sb.WriteString("Type: ")
	sb.WriteString(m.RedisData.GetType())
	sb.WriteString(" ")
	if tm := m.RedisData.GetExpiration(); tm != nil {
		sb.WriteString("Expire: ")
		sb.WriteString(tm.Format(time.RFC3339))
	}
	if len(m.Values) != 0 {
		sb.WriteString("\n")
		sb.WriteString(m.Values.String())
	} else {

	}
	return sb.String()
}
