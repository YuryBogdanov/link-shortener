package handler

import (
	"io"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/config"
	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleNewLinkRegistration() http.HandlerFunc {
	return withLogging(newLinkRegistration())
}

func HandleExistingLinkRequest() http.HandlerFunc {
	return withLogging(existingLinkRequest())
}

func newLinkRegistration() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if url, err := io.ReadAll(r.Body); err == nil {
			linkID, err := storage.MakeAndStoreShortURL(string(url))
			if err != nil {
				handleError(w)
				return
			}
			resultLink := getShortenedLink(r, linkID)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(resultLink))
		} else {
			handleError(w)
		}
	}
	return fn
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

func getShortenedLink(r *http.Request, linkID string) string {
	return config.BaseConfig.ShoretnedBaseURL.Value + "/" + linkID
}

func handleError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}
