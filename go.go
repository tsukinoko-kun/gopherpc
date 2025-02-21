package gopherpc

import (
	"encoding/json"
	"net/http"
	"path"
)

type (
	rpcRequest struct {
		FuncName string `json:"func_name"`
		Args     []any  `json:"args"`
	}

	rpcResponse struct {
		Type   string `json:"type"`
		Result any    `json:"result"`
	}

	rpcError struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}

	muxHandleFunc interface {
		HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	}

	muxHandle interface {
		Handle(pattern string, handler http.Handler)
	}

	muxGet interface {
		Get(pattern string, handler http.HandlerFunc)
	}
)

func handlerFuncRpc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	jr := json.NewDecoder(r.Body)
	req := rpcRequest{}
	if err := jr.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jw := json.NewEncoder(w)
		jw.Encode(rpcError{Error: err.Error(), Type: "error"})
		return
	}

	res, err := call(r, req.FuncName, req.Args)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jw := json.NewEncoder(w)
		jw.Encode(rpcError{Error: err.Error(), Type: "error"})
		return
	} else {
		jw := json.NewEncoder(w)
		jw.Encode(rpcResponse{Result: res, Type: "ok"})
	}
}

func handlerFuncJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	_, _ = w.Write(gopherpcJs)
}

func HandleFunc(mux muxHandleFunc) {
	mux.HandleFunc("/__gopherpc__/rpc", handlerFuncRpc)
	mux.HandleFunc(path.Join("/__gopherpc__", gopherpcJsName), handlerFuncJs)
}

func Handle(mux muxHandle) {
	mux.Handle("/__gopherpc__/rpc", http.HandlerFunc(handlerFuncRpc))
	mux.Handle(path.Join("/__gopherpc__", gopherpcJsName), http.HandlerFunc(handlerFuncJs))
}

func Get(mux muxGet) {
	mux.Get("/__gopherpc__/rpc", handlerFuncRpc)
	mux.Get(path.Join("/__gopherpc__", gopherpcJsName), handlerFuncJs)
}
