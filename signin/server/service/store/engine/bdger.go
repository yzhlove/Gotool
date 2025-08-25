package engine

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/yzhlove/Gotool/signin/server/service/store"
)

const dbpath = "."

type badgerdb struct {
	_db *badger.DB
}

func New() (store.DBer, error) {
	opts := badger.DefaultOptions(dbpath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &badgerdb{db}, err
}

func (s *badgerdb) Set(key, value string) (err error) {
	err = s._db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value))
		return txn.SetEntry(entry)
	})
	return
}

func (s *badgerdb) Get(key string) (value string, err error) {
	err = s._db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			bytes := make([]byte, len(val))
			copy(bytes, val)
			value = string(bytes)
			return nil
		})
	})
	return
}
