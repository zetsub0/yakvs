package api

import (
	"net/http"
	"strings"
)

func checkJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			sendJSONResponse(w, http.StatusUnsupportedMediaType, errResponse{
				Code:  http.StatusUnsupportedMediaType,
				Error: "Content-Type must be application/json",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
