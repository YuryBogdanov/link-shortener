package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/YuryBogdanov/link-shortener/internal/model"
	"github.com/YuryBogdanov/link-shortener/internal/storage"
)

func HandleShortenRequest(w http.ResponseWriter, r *http.Request) {
	if err := validateHeaders(r.Header); err != nil {
		handleError(w)
		return
	}

	var reqModel model.ShortenRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w)
		return
	}
	defer r.Body.Close()

	unmarshalErr := json.Unmarshal(body, &reqModel)
	if unmarshalErr != nil {
		handleError(w)
		return
	}

	link, err := storage.MakeAndStoreShortURL(reqModel.URL)
	if err != nil {
		handleError(w)
		return
	}

	response := prepareResponse(w, link)
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func prepareResponse(w http.ResponseWriter, link string) []byte {
	result := model.Result{Result: link}
	bytes, err := json.Marshal(result)
	if err != nil {
		handleError(w)
	}
	return bytes
}

func validateHeaders(headers http.Header) error {
	contentType := headers.Get("Content-Type")
	if contentType == "application/json" {
		return nil
	} else {
		return errors.New("wrong content type")
	}
}
