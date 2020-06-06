package cache

import (
	"cache/singleflight"
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    GetterFunc
	mainCache cache
	loader *singleflight.Group
}


type GetterFunc func(key string) ([]byte, error)

var (
	mu    sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxBytes int64, getter GetterFunc) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		loader: &singleflight.Group{},
		mainCache: cache{
			maxBytes: maxBytes,
		},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[SimpleCache hit]")
		return v, nil
	}

	view, err := g.loader.Do(key, func() (interface{}, error) {
		log.Printf("[SimpleCache load for database]")
		return g.getLocally(key)
	})

	if err != nil {
		return ByteView{}, err
	}

	return view.(ByteView), nil
}


func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{
		cloneBytes(bytes),
	}
	g.mainCache.add(key, value)
	return value, nil
}
