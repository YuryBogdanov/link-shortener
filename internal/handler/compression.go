package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

const (
	incomingContentCompressionHeader = "Accept-Encoding"
	outgoingContentCompressionHeader = "Content-Encoding"

	encodingMethod = "gzip"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func withCompression(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !needsCompression(r.Header) {
			next.ServeHTTP(w, r)
			return
		}
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			handleError(w)
			return
		}
		gunzippedBody, err := io.ReadAll(gr)
		if err != nil {
			handleError(w)
			return
		}
		buf := bytes.NewBuffer(gunzippedBody)
		r.Body = io.NopCloser(buf)

		w.Header().Add("Content-Encoding", encodingMethod)
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		next.ServeHTTP(w, r)
	}
}

func needsCompression(headers http.Header) bool {
	accept := headers.Get("Accept-Encoding")
	contentEncoding := headers.Get("Content-Encoding")

	return !(contentEncoding == "" || accept == "")
}
