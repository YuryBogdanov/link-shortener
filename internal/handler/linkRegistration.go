package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleNewLinkRegistration(w http.ResponseWriter, r *http.Request) {
	payload, payloadErr := io.ReadAll(r.Body)
	if payloadErr != nil {
		handleError(w)
		return
	}
	url, urlErr := url.ParseRequestURI(string(payload))
	if urlErr != nil {
		handleError(w)
		return
	}
	link, linkErr := storage.MakeAndStoreShortURL(string(url.String()))
	if linkErr != nil {
		handleError(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(link))
}
