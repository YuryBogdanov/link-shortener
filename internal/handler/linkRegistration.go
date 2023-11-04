package handler

import (
	"io"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/config"
	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleNewLinkRegistration(w http.ResponseWriter, r *http.Request) {
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

func getShortenedLink(r *http.Request, linkID string) string {
	return config.BaseConfig.ShoretnedBaseURL.Value + "/" + linkID
}
