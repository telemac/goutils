package main

import (
	//log "github.com/sirupsen/logrus"
	"log"

	"github.com/telemac/goutils/keyvalue"
)

func main() {
	log.Println("kv test")

	// Create a kv store
	kv := keyvalue.NewKVMapStore()
	kv.Set(keyvalue.Key("name"), keyvalue.Value("Alexandre"))
	name, err := kv.Get(keyvalue.Key("name"))
	if err == keyvalue.ErrKeyNotFound {
		log.Fatalf("key %s not found", "name")
	}
	log.Printf("name is %s", name)

}
