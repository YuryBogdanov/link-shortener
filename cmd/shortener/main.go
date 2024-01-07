package main

import (
	"log"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/config"
	"github.com/YuryBogdanov/link-shortener/internal/handler"
	"github.com/YuryBogdanov/link-shortener/internal/logger"
	"github.com/YuryBogdanov/link-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

var lg = logger.DefaultLogger{}

func main() {
	config.SetupFlags()

	lg.Setup()
	defer lg.Finish()

	router := chi.NewRouter()
	router.Get("/{id}", handler.HandleExistingLinkRequest())
	router.Post("/", handler.HandleNewLinkRegistration())
	router.Post("/api/shorten", handler.HandleShortenRequest())

	baseURL := config.BaseConfig.ServerPath.Value
	lg.Info(
		"Starting server",
		"addr", baseURL,
	)

	storageFile := config.BaseConfig.StorageFilePath.Value
	storage.SetupPersistentStorage(storageFile)

	err := http.ListenAndServe(baseURL, router)
	if err != nil {
		log.Fatal(err)
	}
}
