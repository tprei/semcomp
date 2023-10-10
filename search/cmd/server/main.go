package main

import (
	"log"

	"github.com/tprei/semcomp/search/internal/server"
)

func main() {
	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
