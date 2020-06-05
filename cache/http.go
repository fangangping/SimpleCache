package cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_simplecache/"

type HTTPPool struct {
	self string
	basePath string
}


func NewHttpPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, "wrong base path", http.StatusBadRequest)
		return
	}

	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "wrong path", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	keyName := parts[1]

	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, fmt.Sprintf("no such group: %s", groupName), http.StatusBadRequest)
		return
	}

	view, err := g.Get(keyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}