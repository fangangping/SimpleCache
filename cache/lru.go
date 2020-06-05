package cache

import (
	"container/list"
)

type Lru struct {
	maxBytes  int64
	currBytes int64
	list      *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key string
	value Value
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Lru {
	return &Lru{
		maxBytes:  maxBytes,
		currBytes: 0,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (lc *Lru) Get(key string) (value Value, ok bool){
	if ele, ok := lc.cache[key]; ok {
		lc.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (lc *Lru) Add(key string, value Value) {
	if ele, ok := lc.cache[key]; ok {
		lc.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		lc.currBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := lc.list.PushFront(&entry{
			key: key,
			value: value,
		})
		lc.cache[key] = ele
		lc.currBytes += int64(len(key)) + int64(value.Len())
	}

	for lc.maxBytes != 0 && lc.currBytes > lc.maxBytes {
		lc.removeOldest()
	}
}

func (lc *Lru) removeOldest() {
	ele := lc.list.Back()
	if ele != nil {
		lc.list.Remove(ele)
		kv := ele.Value.(*entry)
		delete(lc.cache, kv.key)
		lc.currBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if lc.OnEvicted != nil {
			lc.OnEvicted(kv.key, kv.value)
		}
	}
}

func (lc *Lru) Len() int {
	return lc.list.Len()
}
