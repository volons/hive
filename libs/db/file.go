package db

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger"
)

// FileDB is a simple bolt based persistent key value store
type FileDB struct {
	path string
	data *badger.DB
}

// NewFileDB creates a FileDB given a path to store data on the filesystem
func NewFileDB(path string) *FileDB {
	return &FileDB{
		path: path,
	}
}

// Get returns the value for the given key
func (db *FileDB) Get(key string, valPtr interface{}) error {
	return db.data.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		val, err := item.Value()
		if err != nil {
			return err
		}

		return json.Unmarshal(val, valPtr)
	})
}

// Set sets the value of the given key
func (db *FileDB) Set(key string, val interface{}) error {
	value, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return db.data.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		return err
	})
}

// Delete removes the key from the database
func (db *FileDB) Delete(key string) error {
	return db.data.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// SetWithTTL sets the value of the given key with an time to live (expiration)
func (db *FileDB) SetWithTTL(key string, val interface{}, d time.Duration) error {
	value, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return db.data.Update(func(txn *badger.Txn) error {
		err := txn.SetWithTTL([]byte(key), value, d)
		return err
	})
}

func (db *FileDB) Find(prefix string) ([]string, error) {
	var keys []string

	err := db.data.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
		}

		return nil
	})

	return keys, err
}

// Init opens the database
func (db *FileDB) Init() error {
	opts := badger.DefaultOptions
	opts.Dir = db.path
	opts.ValueDir = db.path
	data, err := badger.Open(opts)
	db.data = data
	return err
}

func IsNotFoudError(err error) bool {
	return err == badger.ErrKeyNotFound
}
