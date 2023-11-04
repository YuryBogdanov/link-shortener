package main

import (
	"log"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/config"
	"github.com/YuryBogdanov/link-shortener/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.SetupFlags()

	router := chi.NewRouter()
	router.Get("/{id}", handler.HandleExistingLinkRequest)
	router.Post("/", handler.HandleNewLinkRegistration)
	router.Post("/api/shorten", handler.HandleShortenRequest)

	baseURL := config.BaseConfig.ServerPath.Value
	err := http.ListenAndServe(baseURL, router)
	if err != nil {
		log.Fatal(err)
	}
}
