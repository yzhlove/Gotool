package rdb

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"time"

	"github.com/gookit/goutil"
	"github.com/hdt3213/rdb/parser"
	"rain.com/Gotool/redis-dump/app/internal/decode"
)

func Dump(reader io.Reader) ([]Meta, error) {
	var values []Meta
	var dec = parser.NewDecoder(reader)
	var err error
	err = dec.Parse(func(o parser.RedisObject) bool {
		switch o.GetType() {
		case parser.StringType:
			var m = Meta{RedisData: o}
			str := o.(*parser.StringObject)
			if ret, value := decode.Parse(str.Value); ret != decode.Unknown {
				m.Values = append(m.Values, Data{Type: ret, Value: value})
			} else {
				m.Values = append(m.Values, Data{Type: ret, Value: goutil.String(str.Value)})
			}
			values = append(values, m)
		case parser.ListType:
			list := o.(*parser.ListObject)
			var m = Meta{RedisData: o}
			for _, bytes := range list.Values {
				if ret, value := decode.Parse(bytes); ret != decode.Unknown {
					m.Values = append(m.Values, Data{Type: ret, Value: value})
				} else {
					m.Values = append(m.Values, Data{Type: ret, Value: goutil.String(bytes)})
				}
			}
			values = append(values, m)
		case parser.HashType:
			hash := o.(*parser.HashObject)
			var m = Meta{RedisData: o}
			var keys = slices.Collect(maps.Keys(hash.Hash))
			slices.Sort(keys)
			for _, key := range keys {
				var extra = key
				var bytes = hash.Hash[key]
				if expire := hash.FieldExpirations[key]; expire != 0 {
					t := time.Unix(expire, 0)
					extra = key + ":" + t.Format(time.RFC3339)
				}
				if ret, value := decode.Parse(bytes); ret != decode.Unknown {
					m.Values = append(m.Values, Data{Name: extra, Type: ret, Value: value})
				} else {
					m.Values = append(m.Values, Data{Name: extra, Type: ret, Value: goutil.String(bytes)})
				}
			}
			values = append(values, m)
		case parser.ZSetType:
			zset := o.(*parser.ZSetObject)
			var m = Meta{RedisData: o}
			for _, e := range zset.Entries {
				if ret, value := decode.Parse([]byte(e.Member)); ret != decode.Unknown {
					m.Values = append(m.Values,
						Data{Name: fmt.Sprintf("score:%#v", e.Score), Type: ret, Value: value})
				} else {
					m.Values = append(m.Values,
						Data{Name: fmt.Sprintf("score:%#v", e.Score),
							Type: ret, Value: goutil.String(e.Member)})
				}
			}
			values = append(values, m)
		case parser.SetType:
			set := o.(*parser.SetObject)
			var m = Meta{RedisData: o}
			for _, bytes := range set.Members {
				if ret, value := decode.Parse(bytes); ret != decode.Unknown {
					m.Values = append(m.Values, Data{Type: ret, Value: value})
				} else {
					m.Values = append(m.Values, Data{Type: ret, Value: goutil.String(bytes)})
				}
			}
			values = append(values, m)
		case parser.AuxType:
			aux := o.(*parser.AuxObject)
			var m = Meta{RedisData: o}
			m.Values = append(m.Values, Data{Type: decode.Unknown, Value: goutil.String(aux)})
			values = append(values, m)
		case parser.DBSizeType:
			dbsize := o.(*parser.DBSizeObject)
			var m = Meta{RedisData: o}
			m.Values = append(m.Values, Data{Type: decode.Unknown, Value: goutil.String(dbsize)})
			values = append(values, m)
		}
		return true
	})
	return values, err
}
