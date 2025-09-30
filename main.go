package main

import (
	"log"

	httpController "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/Arup3201/gotasks/internal/storages"
	. "github.com/Arup3201/gotasks/internal/utils"
)

func main() {
	Config.Configure()

	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("Storage creation failed: %v", err)
	}

	err = httpController.InitServer(storage)
	if err != nil {
		log.Fatalf("Server create failed: %v", err)
	}
	err = httpController.Server.Run("0.0.0.0")
	if err != nil {
		log.Fatalf("Server create failed: %v", err)
	}
}
