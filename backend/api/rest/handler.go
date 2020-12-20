package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRootHandler() RootHandler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(AttachRequestIDInLogger)

	router.Get("/health", health)

	return RootHandler{
		router: router,
	}
}

type RootHandler struct {
	router *chi.Mux
}

func (h RootHandler) Get() http.Handler {
	return h.router
}

type CustomHandler interface {
	register(router *chi.Mux)
}

func (h RootHandler) Register(handlers ...CustomHandler) {
	for index := range handlers {
		handlers[index].register(h.router)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}