package main

import (
	"log"

	"github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/Arup3201/gotasks/internal/storages"
	. "github.com/Arup3201/gotasks/internal/utils"
)

func main() {
	Config.Configure()

	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("Storage creation failed: %v", err)
	}

	server := http.CreateServer(storage)
	server.AttachRoutes()
	err = server.Run("localhost") // default localhost, otherwise pass the URL
	if err != nil {
		log.Fatalf("Server create failed: %v", err)
	}
}
