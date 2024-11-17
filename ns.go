package gopherpc

import (
	"context"
	"fmt"
)

type AnyFunc func(context.Context, []any) (any, error)

var ns = map[string]AnyFunc{}

func Register(name string, f AnyFunc) {
	ns[name] = f
}

func call(ctx context.Context, name string, args []any) (res any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("function %q panic: %v", name, r)
			res = nil
		}
	}()

	if f, ok := ns[name]; ok {
		res, err = f(ctx, args)
		return
	}

	return nil, fmt.Errorf("function %q not found", name)
}