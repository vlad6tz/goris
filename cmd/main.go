package main

import (
	"fmt"

	"github.com/vlad6tz/goris.git/internal/server"
	"github.com/vlad6tz/goris.git/internal/store"
)

func main() {
	db := store.New()
	serv, err := server.New("localhost:8080", db)
	if err != nil {
		fmt.Printf("error creating new server: ", err)
		return
	}

	serv.Serve()
}
