package handler

import (
	"net/http"
	"time"

	"github.com/YuryBogdanov/link-shortener/internal/logger"
)

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Info(
			"Serving Request",
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}

	return http.HandlerFunc(logFunc)
}
