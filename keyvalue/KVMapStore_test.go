package keyvalue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testKVStore(t *testing.T, kv KVStore) {
	assert := assert.New(t)

	kvType := fmt.Sprintf("keyvlaue store type is %T", kv)

	err := kv.Set(Key("name"), Value("Alexandre"))
	assert.Nil(err, kvType)

	v, err := kv.Get(Key("name"))
	assert.NoError(err)
	assert.Equal(Value("Alexandre"), v)

	_, err = kv.Get(Key("unknown_key"))
	assert.Equal(ErrKeyNotFound, err, kvType)
}

func TestKVMapStore(t *testing.T) {

	kv := NewKVMapStore()
	testKVStore(t, kv)

}
