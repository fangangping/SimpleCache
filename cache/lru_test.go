package cache

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestLruCache_Get(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Add("key1", String("value1"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "value1" {
		t.Fatal("cache hit key != 1234 failed")
	}

	lru.Add("key1", String("value2"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "value2" {
		t.Fatal("cache hit key != 1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

func TestLruCache_removeOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1) + len(v1) + len(k2) + len(v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatal("RemoveOldest key1 failed")
	}
}

func TestLruCache_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}


	cap := len("key3") + len("value3") + len("key4") + len("value4")

	lru := New(int64(cap), callback)

	lru.Add("key1", String("value1"))
	lru.Add("key2", String("value2"))
	lru.Add("key3", String("value3"))
	lru.Add("key4", String("value3"))

	expect := []string{"key1", "key2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s, get %s", expect, keys)
	}

}