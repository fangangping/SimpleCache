package cache

import (
	"fmt"
	"sync"
	"testing"
)

var db = map[string]string{
	"key1": "value1",
	"key2": "value2",
	"key3": "value3",
}

func TestCache(t *testing.T) {
	loadCount := make(map[string]int)

	callback := func(key string) ([]byte, error){
		if v, ok := db[key]; ok {
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			}
			loadCount[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %s not exis", key)
	}
	g := NewGroup("test", 0, callback)

	if g != GetGroup("test") {
		t.Fatalf("fail to get gruop")
	}

	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatalf("expect to get %s, but get %s", v, view.String())
		}
		if view, err := g.Get(k); err != nil || view.String() != v || loadCount[k] != 1 {
			t.Fatalf("%s cache miss", k)
		}
	}

	if v, err := g.Get(""); err == nil {
		t.Fatalf("empty key should get err but get %s", v)
	}

	if v, err := g.Get("unknown"); err == nil {
		t.Fatalf("unknown key should get err but get %s", v)
	}
}

func TestSingleFlight(t *testing.T) {
	loadCount := make(map[string]int)

	mu := sync.Mutex{}
	callback := func(key string) ([]byte, error){
		mu.Lock()
		defer mu.Unlock()
		if v, ok := db[key]; ok {
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			}
			loadCount[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %s not exis", key)
	}
	g := NewGroup("test", 0, callback)

	if g != GetGroup("test") {
		t.Fatalf("fail to get gruop")
	}

	for i := 0; i < 100; i++ {
		go func() {
			for k, v := range db {
				if view, err := g.Get(k); err != nil || view.String() != v {
					t.Fatalf("expect to get %s, but get %s", v, view.String())
				}
				if view, err := g.Get(k); err != nil || view.String() != v || loadCount[k] != 1 {
					t.Fatalf("%s cache miss", k)
				}
			}
		}()
	}


	if v, err := g.Get(""); err == nil {
		t.Fatalf("empty key should get err but get %s", v)
	}

	if v, err := g.Get("unknown"); err == nil {
		t.Fatalf("unknown key should get err but get %s", v)
	}
}
