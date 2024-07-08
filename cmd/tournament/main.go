package main

import (
	"github.com/anrew1002/Tournament-ChemLoto/internal/server"
)

func main() {
	port := "1090"
	server := server.NewServer()
	server.Run(port)
}
