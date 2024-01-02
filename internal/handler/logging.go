package handler

import (
	"net/http"
	"time"

	"github.com/YuryBogdanov/link-shortener/internal/logger"
)

var lg = logger.DefaultLogger{}

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	lg.Setup()
	defer lg.Finish()

	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)

		duration := time.Since(start)

		lg.Info(
			"Serving Request",
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}

	return http.HandlerFunc(logFunc)
}
