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
		if link, err := storage.GetLinkForKey(query); err == nil {
			w.Header().Add("Location", link)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			handleError(w)
		}
	}
	return fn
}
