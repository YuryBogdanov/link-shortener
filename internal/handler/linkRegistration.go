package handler

import (
	"io"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleNewLinkRegistration() http.HandlerFunc {
	return withLogging(newLinkRegistration())
}

func newLinkRegistration() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if url, err := io.ReadAll(r.Body); err == nil {
			linkID, err := storage.MakeAndStoreShortURL(string(url))
			if err != nil {
				handleError(w)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(linkID))
		} else {
			handleError(w)
		}
	}
	return fn
}
