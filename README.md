# GopheRPC

GopheRPC is a RPC library that lets you call Go server functions from your JavaScript (browser) clients.

```go
package main

import (
	"context"
	"net/http"

	"github.com/tsukinoko-kun/gopherpc"
)

func main() {
	gopherpc.Register("foo", func(ctx context.Context, args []any) (any, error) {
		return "bar", nil
	})
	gopherpc.GopheRPC(http.DefaultServeMux)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><head><title>GopherPC</title></head><body>` +
			gopherpc.ImportJs() +
			`<button onclick="gopherpc.foo().then(alert)">Call foo</button></body></html>`))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
```
