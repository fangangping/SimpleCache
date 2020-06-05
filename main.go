package main

import (
	"SimpleCache/cache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"key1": "value1",
	"key2": "value2",
	"key3": "value3",
}

func main() {
	cache.NewGroup("test", 2 << 10, func(key string) ([]byte, error){
		log.Printf("[DB] search key %s", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %s does not exist", key)
	})

	addr := "0.0.0.0:8080"
	server := cache.NewHttpPool(addr)
	log.Printf("server is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}