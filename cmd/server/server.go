package main

import (
	"fmt"
	"github.com/felixge/felixge.de/fs"
	"net"
	"net/http"
)

func main() {
	fs := fs.New()

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening at: http://%s\n", listener.Addr())

	if err := http.Serve(listener, http.FileServer(fs)); err != nil {
		panic(err)
	}
}
