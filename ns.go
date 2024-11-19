package gopherpc

import (
	"fmt"
	"net/http"
	"sync"
)

type AnyFunc func(r *http.Request, args []any) (any, error)

var (
	ns  = map[string]AnyFunc{}
	mut = sync.RWMutex{}
)

func Register(name string, f AnyFunc) {
	mut.Lock()
	defer mut.Unlock()
	ns[name] = f
}

func Remove(name string) {
	mut.Lock()
	defer mut.Unlock()
	delete(ns, name)
}

func call(r *http.Request, name string, args []any) (res any, err error) {
	mut.RLock()
	defer mut.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("function %q panic: %v", name, r)
			res = nil
		}
	}()

	if f, ok := ns[name]; ok {
		res, err = f(r, args)
		return
	}

	return nil, fmt.Errorf("function %q not found", name)
}
