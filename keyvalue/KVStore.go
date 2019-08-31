package keyvalue

import "errors"

var (
	errKeyEmpty    = errors.New("key is empty")
	errKeyNotFound = errors.New("key not found")
)

// Key is the type for a key
type Key []byte

// Value is the type for a value
type Value []byte

// KVStore is a simple interface for KeyValue stores
type KVStore interface {
	// Set sets or replaces a key/value pair in the store
	Set(key Key, value Value) error

	// Get gets the value given the key
	Get(key Key) (Value, error)

	// Delete deletes a key, returns an error if the key don't exist
	Delete(key Key) error
}
