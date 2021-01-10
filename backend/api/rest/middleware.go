package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"

	"github.com/charlesvdv/cirrus/backend/pkg/identity"
)

type contextKey int

const (
	ctxKeyMarshaller contextKey = iota
	ctxKeyUser
)

func getUserID(ctx context.Context) identity.UserID {
	return ctx.Value(ctxKeyUser).(identity.UserID)
}

type CheckBearerTokenFn func(context.Context, string) (identity.UserID, error)

func authMiddleware(fn CheckBearerTokenFn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Ctx(r.Context()).Warn().Msg("auth middleware")
			token := getAuthenticationToken(r)
			userID, err := fn(r.Context(), token)
			if err != nil {
				renderError(r.Context(), w, convertIdentityErr(err))
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUser, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getAuthenticationToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func marshallingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var contentType string
		if r.Body == nil || r.ContentLength == 0 {
			contentType = "application/json"
		} else {
			contentType = r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "", http.StatusUnsupportedMediaType)
				return
			}
		}

		w.Header().Add("Content-Type", contentType)

		ctx := context.WithValue(r.Context(), ctxKeyMarshaller, jsonMarshaller{})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func attachRequestIDInLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("request-id", middleware.GetReqID(r.Context())).Logger()
		ctx := logger.WithContext(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
