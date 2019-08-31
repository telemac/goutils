package keyvalue

// KVMapStore implements a non thread safe KVStore based on go maps
type KVMapStore struct {
	store map[string]Value
}

// NewKVMapStore creates a KVMapStore
func NewKVMapStore() *KVMapStore {
	var kvMapStore KVMapStore
	kvMapStore.store = make(map[string]Value)
	return &kvMapStore
}

// Set adds or replaces the key / value pair in the map
func (kvms *KVMapStore) Set(key Key, value Value) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	kvms.store[string(key)] = value
	return nil
}

// Get returns the value associated with the key, or an error if the key does not exist
func (kvms *KVMapStore) Get(key Key) (Value, error) {
	value, found := kvms.store[string(key)]
	if !found {
		return Value(""), errKeyNotFound
	}
	return value, nil
}

// Delete deletes the key from the map, or returns an error if the key is not found
func (kvms KVMapStore) Delete(key Key) error {
	_, found := kvms.store[string(key)]
	if !found {
		return errKeyNotFound
	}
	delete(kvms.store, string(key))
	return nil
}
