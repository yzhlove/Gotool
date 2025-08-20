package main

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
)

func main() {

	defaultOptions := badger.DefaultOptions("./badger/chat07")
	defaultOptions.Logger = nil

	db, err := badger.Open(defaultOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("language"), []byte("golang"))
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte("peoples"), []byte("china people!"))
		return txn.SetEntry(entry)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("language"))
		if err != nil {
			return err
		}

		if err = item.Value(func(val []byte) error {
			fmt.Println(string(val))
			return nil
		}); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100

		iterations := txn.NewIterator(opts)
		defer iterations.Close()

		for iterations.Rewind(); iterations.Valid(); iterations.Next() {
			item := iterations.Item()
			key := item.Key()
			err = item.Value(func(val []byte) error {
				fmt.Println(string(key), string(val))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

}
