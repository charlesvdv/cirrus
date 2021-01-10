package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRootHandler() *RootHandler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(attachRequestIDInLogger)
	router.Use(middleware.Recoverer)

	router.Get("/health", health)

	return &RootHandler{
		router: router,
	}
}

type handlerContext struct {
	authMiddleware func(http.Handler) http.Handler
}

type RootHandler struct {
	router     *chi.Mux
	handlerCtx handlerContext
}

func (h RootHandler) Get() http.Handler {
	return h.router
}

type CustomHandler interface {
	register(router *chi.Mux, ctx handlerContext)
}

func (h *RootHandler) WithTokenBearerChecker(fn CheckBearerTokenFn) *RootHandler {
	h.handlerCtx.authMiddleware = authMiddleware(fn)
	return h
}

func (h *RootHandler) Register(handlers ...CustomHandler) {
	for index := range handlers {
		handlers[index].register(h.router, h.handlerCtx)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
