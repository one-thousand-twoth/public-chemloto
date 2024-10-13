package main

import (
	"os"

	"github.com/anrew1002/Tournament-ChemLoto/internal/server"
)

func main() {
	initResultFolder()
	port := "8090"
	server := server.NewServer()
	server.Run(port)
}
func initResultFolder() {
	err := os.Mkdir("chemloto results", 0755)
	if err != nil {
		// Если директория уже существует
		if !os.IsExist(err) {
			panic(
				"Cannot create folder for results!",
			)
		}
		return
	}
}
