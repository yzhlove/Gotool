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

func Dump(reader io.Reader, callback func(meta *Meta)) error {
	var dec = parser.NewDecoder(reader)
	var err error
	err = dec.Parse(func(o parser.RedisObject) bool {
		var m = &Meta{RedisData: o}
		switch o.GetType() {
		case parser.StringType:
			str := o.(*parser.StringObject)
			if ret, value := decode.Parse(str.Value); ret != decode.Unknown {
				m.Values = append(m.Values, Data{Type: ret, Value: value})
			} else {
				m.Values = append(m.Values, Data{Type: ret, Value: goutil.String(str.Value)})
			}
			if callback != nil {
				callback(m)
			}
		case parser.ListType:
			list := o.(*parser.ListObject)
			for _, bytes := range list.Values {
				if ret, value := decode.Parse(bytes); ret != decode.Unknown {
					m.Values = append(m.Values, Data{Type: ret, Value: value})
				} else {
					m.Values = append(m.Values, Data{Type: ret, Value: goutil.String(bytes)})
				}
			}
			if callback != nil {
				callback(m)
			}
		case parser.HashType:
			hash := o.(*parser.HashObject)
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
			if callback != nil {
				callback(m)
			}
		case parser.ZSetType:
			zset := o.(*parser.ZSetObject)
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
			if callback != nil {
				callback(m)
			}
		case parser.SetType:
			set := o.(*parser.SetObject)
			for _, bytes := range set.Members {
				if ret, value := decode.Parse(bytes); ret != decode.Unknown {
					m.Values = append(m.Values, Data{Type: ret, Value: value})
				} else {
					m.Values = append(m.Values, Data{Type: ret, Value: goutil.String(bytes)})
				}
			}
			if callback != nil {
				callback(m)
			}
		case parser.AuxType:
			aux := o.(*parser.AuxObject)
			m.Values = append(m.Values, Data{Type: decode.Unknown, Value: goutil.String(aux)})
			if callback != nil {
				callback(m)
			}
		case parser.DBSizeType:
			dbsize := o.(*parser.DBSizeObject)
			m.Values = append(m.Values, Data{Type: decode.Unknown, Value: goutil.String(dbsize)})
			if callback != nil {
				callback(m)
			}
		}
		return true
	})
	return err
}
