package main

import (
	"io"
	"net/http"
)

func processShortURLRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleNewLinkRegistration(w, r)

	case http.MethodGet:
		handleExistingLinkRequest(w, r)

	default:
		handleError(w)
	}
}

func handleNewLinkRegistration(w http.ResponseWriter, r *http.Request) {
	if url, err := io.ReadAll(r.Body); err == nil {
		linkID, err := makeAndStoreShortURL(string(url))
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

func handleExistingLinkRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path[1:]
	if link, exists := links[query]; exists {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		handleError(w)
	}
}

func getShortenedLink(r *http.Request, linkID string) string {
	return "http://" + r.Host + "/" + linkID
}

func handleError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", processShortURLRequest)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
