package gopherpc

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

type AnyFunc func(r *http.Request, args []any) (any, error)

var (
	ns     = map[string]AnyFunc{}
	nsTyps = map[string]string{}
	mut    = sync.RWMutex{}
)

func Types() string {
	mut.RLock()
	defer mut.RUnlock()

	ts := strings.Builder{}
	ts.WriteString("export {};\n")
	ts.WriteString("declare global {\n")
	ts.WriteString("  const gopherpc: {\n")
	for name := range ns {
		if t, ok := nsTyps[name]; ok {
			ts.WriteString(fmt.Sprintf("    %q: %s;\n", name, t))
		} else {
			ts.WriteString(fmt.Sprintf("    %q: <R>(...args: any) => Promise<R>;\n", name))
		}
	}
	ts.WriteString("  };\n")
	ts.WriteString("}\n")

	return ts.String()
}

func RegisterTyped[A any, R any](name string, f func(r *http.Request, arg A) (R, error)) {
	mut.Lock()
	defer mut.Unlock()

	fType := reflect.TypeOf(f)

	var arguments []string
	argType := fType.In(1)
	if argType.Kind() == reflect.Struct {
		numFields := argType.NumField()
		for i := 0; i < numFields; i++ {
			field := argType.Field(i)
			arguments = append(arguments, fmt.Sprintf("%s: %s", field.Name, goTypeToTypescript(field.Type)))
		}
	} else {
		arguments = append(arguments, fmt.Sprintf("arg: %s", goTypeToTypescript(argType)))
	}

	ns[name] = func(r *http.Request, args []any) (any, error) {
		a := reflect.New(argType).Elem().Interface().(A)

		if err := Unmarshal(args, &a); err != nil {
			return nil, err
		}

		log.Printf("%s(%v) : %v\n", name, args, a)

		return f(r, a)
	}
	nsTyps[name] = fmt.Sprintf("(%s) => Promise<%s>", strings.Join(arguments, ", "), goTypeToTypescript(fType.Out(0)))
}

func Register(name string, f AnyFunc) {
	mut.Lock()
	defer mut.Unlock()
	ns[name] = f
}

func Remove(name string) {
	mut.Lock()
	defer mut.Unlock()
	delete(ns, name)
	delete(nsTyps, name)
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
