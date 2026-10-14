package main

import (
	"hw1/internal/cli"
	"hw1/internal/storage"
	"hw1/internal/app"
)

func main() {
	storage := storage.MakeNewStorageMap()
	
	handler := app.MakeNewTaskHandler(storage)

	cli.Run(handler)
}
