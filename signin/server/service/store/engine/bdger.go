package engine

import (
	"github.com/dgraph-io/badger/v4"
)

const dbpath = "."

type store struct {
	_db *badger.DB
}

func New() (*store, error) {
	opts := badger.DefaultOptions(dbpath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &store{db}, err
}

func (s *store) Set(key, value string) (err error) {
	err = s._db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value))
		return txn.SetEntry(entry)
	})
	return
}

func (s *store) Get(key string) (value string, err error) {
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
