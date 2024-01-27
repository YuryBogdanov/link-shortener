package handler

import (
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleExistingLinkRequest() http.HandlerFunc {
	return withLogging(existingLinkRequest())
}

func existingLinkRequest() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Path[1:]
		link, err := storage.GetLinkForKey(query)
		if err != nil {
			handleError(w)
			return
		}
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
	return fn
}
