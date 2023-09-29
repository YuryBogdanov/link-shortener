package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
)

var links map[string]string

func processShortURLRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		} else {
			linkId := makeAndStoreShortURL(string(url))
			resultLink := getShortenedLink(r, linkId)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(resultLink))
		}
	} else if r.Method == http.MethodGet {
		query := r.URL.Path[1:]
		if link, exists := links[query]; exists {
			w.Header().Add("Location", link)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			handleError(w)
		}
	} else {
		handleError(w)
	}
}

func makeAndStoreShortURL(url string) string {
	hash := md5.New()
	io.WriteString(hash, url)
	encodedString := fmt.Sprintf("%x", hash.Sum(nil))
	if len([]rune(encodedString)) < 8 {
		links[encodedString] = url
		return encodedString
	} else {
		links[encodedString[:8]] = url
		return encodedString[:8]
	}
}

func getShortenedLink(r http.Request, linkId string) string {
	return "http://" + r.Host + "/" + linkId
}

func handleError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	links = make(map[string]string)
	mux := http.NewServeMux()

	mux.HandleFunc("/", processShortURLRequest)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
