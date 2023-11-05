package main

import (
	"log"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/config"
	"github.com/YuryBogdanov/link-shortener/internal/handler"
	"github.com/YuryBogdanov/link-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.SetupFlags()

	logger.Setup()
	defer logger.Finish()

	router := chi.NewRouter()
	router.Get("/{id}", handler.HandleExistingLinkRequest())
	router.Post("/", handler.HandleNewLinkRegistration())
	router.Post("/api/shorten", handler.HandleShortenRequest())

	baseURL := config.BaseConfig.ServerPath.Value
	logger.Info(
		"Starting server",
		"addr", baseURL,
	)
	err := http.ListenAndServe(baseURL, router)
	if err != nil {
		log.Fatal(err)
	}
}
