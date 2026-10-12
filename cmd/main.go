package main

import (
	"fmt"

	"github.com/vlad6tz/goris.git/internal/server"
	"github.com/vlad6tz/goris.git/internal/store"
)

func main() {
	db := store.New()
	serv, err := server.New("localhost:8080", db, 64, 256)
	if err != nil {
		fmt.Printf("error creating new server: %s", err)
		return
	}

	serv.Serve()
}
