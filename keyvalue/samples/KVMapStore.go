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
	err := kv.Set(keyvalue.Key("name"), keyvalue.Value("Alexandre"))
	if err != nil {
		log.Fatalf("Unable to set key %s", "name")
	}
	name, err := kv.Get(keyvalue.Key("name"))
	if err == keyvalue.ErrKeyNotFound {
		log.Fatalf("key %s not found", "name")
	}
	log.Printf("name is %s", name)
}
