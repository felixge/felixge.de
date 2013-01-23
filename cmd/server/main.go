package main

import (
	"github.com/felixge/felixge.de/fs"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	log.Printf("Listening at: http://%s\n", listener.Addr())

	h := newHandler()
	if err := http.Serve(listener, h); err != nil {
		panic(err)
	}
}

func newHandler() *handler {
	fs := fs.New()
	return &handler{fileServer: http.FileServer(fs)}
}

type handler struct {
	fileServer http.Handler
}

func (h *handler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	start := time.Now()
	h.fileServer.ServeHTTP(res, req)
	duration := time.Since(start)

	log.Printf("%s %s took %s", req.Method, req.URL.Path, duration)
}
