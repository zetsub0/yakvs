package api

import (
	"log/slog"
	"net/http"
)

func logReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"got request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
