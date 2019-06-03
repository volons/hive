package db

import "time"

// Database is an interface for a simple key value store
// intended to be backed by the appropriate database engine
// for the intended use case
// ex:
//   - redis if running on a server
//   - bolt if running in an embeded environment (mobile, drone, ...)
type Database interface {
	Init() error
	Get(string, interface{}) error
	Find(string) ([]string, error)
	Set(string, interface{}) error
	SetWithTTL(string, interface{}, time.Duration) error
	Delete(string) error
}

// DB holds the global database instance
var DB Database

// Init is a global function initializes the configured database
func Init(conf string) error {
	DB = NewFileDB(conf)
	return DB.Init()
}

// Get returns the value for the given key
func Get(key string, valPtr interface{}) error {
	return DB.Get(key, valPtr)
}

// Find retruns all keys that start with the specified prefix
func Find(prefix string) ([]string, error) {
	return DB.Find(prefix)
}

// Set sets the value of the given key
func Set(key string, val interface{}) error {
	return DB.Set(key, val)
}

// SetWithTTL sets the value of the given key with an time to live (expiration)
func SetWithTTL(key string, val interface{}, d time.Duration) error {
	return DB.SetWithTTL(key, val, d)
}

// Delete deletes a key from the database
func Delete(key string) error {
	return DB.Delete(key)
}
