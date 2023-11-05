package handler

import (
	"io"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleNewLinkRegistration(w http.ResponseWriter, r *http.Request) {
	if url, err := io.ReadAll(r.Body); err == nil {
		link, err := storage.MakeAndStoreShortURL(string(url))
		if err != nil {
			handleError(w)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(link))
	} else {
		handleError(w)
	}
}
