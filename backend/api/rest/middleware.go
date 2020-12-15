package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

type contextKey int

const (
	ctxKeyMarshaller contextKey = iota
)

func marshallingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "", http.StatusUnsupportedMediaType)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyMarshaller, jsonMarshaller{})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AttachRequestIDInLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("request-id", middleware.GetReqID(r.Context())).Logger()
		ctx := logger.WithContext(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
