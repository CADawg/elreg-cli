package main

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
)

func GetReadStatusBadger(homepageUrl string) bool {
	var localPostInfo LocalPostInfo

	_ = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("v1_" + homepageUrl))

		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &localPostInfo)
		})
	})

	return localPostInfo.IsRead
}

func SetReadStatusBadger(homepageUrl string, readStatus bool) {
	localPostInfo := LocalPostInfo{Url: homepageUrl, IsRead: readStatus}

	_ = db.Update(func(txn *badger.Txn) error {
		encoded, err := json.Marshal(localPostInfo)

		if err != nil {
			return err
		}

		return txn.Set([]byte("v1_"+homepageUrl), encoded)
	})
}
