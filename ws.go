package gopherpc

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type (
	rpcRequest struct {
		FuncName string `json:"func_name"`
		Args     []any  `json:"args"`
		Id       string `json:"id"`
	}

	rpcResponse struct {
		Result any    `json:"result"`
		Id     string `json:"id"`
	}

	rpcError struct {
		Error string `json:"error"`
		Id    string `json:"id"`
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var req rpcRequest
		if err = json.Unmarshal(message, &req); err != nil {
			break
		}

		res, err := call(r, req.FuncName, req.Args)
		if err != nil {
			_ = conn.WriteJSON(rpcError{Error: err.Error(), Id: req.Id})
			continue
		} else {
			_ = conn.WriteJSON(rpcResponse{Result: res, Id: req.Id})
		}
	}
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	_, _ = w.Write(gopherpcJs)
}

func HandleFunc(mux muxHandleFunc) {
	mux.HandleFunc("/__gopherpc__/ws", wsHandler)
	mux.HandleFunc(path.Join("/__gopherpc__", gopherpcJsName), jsHandler)
}

func Handle(mux muxHandle) {
	mux.Handle("/__gopherpc__/ws", http.HandlerFunc(wsHandler))
	mux.Handle(path.Join("/__gopherpc__", gopherpcJsName), http.HandlerFunc(jsHandler))
}

func Get(mux muxGet) {
	mux.Get("/__gopherpc__/ws", wsHandler)
	mux.Get(path.Join("/__gopherpc__", gopherpcJsName), jsHandler)
}
