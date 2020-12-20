package rest

import (
	"context"
	"net/http"

	"github.com/charlesvdv/cirrus/backend/pkg/user"
	"github.com/go-chi/chi"
)

type UserService interface {
	Signup(ctx context.Context, info user.SignupInfo) error
}

func NewUserHandler(service UserService) UserHandler {
	return UserHandler{
		service: service,
	}
}

type UserHandler struct {
	service UserService
}

func (h UserHandler) register(router *chi.Mux) {
	router.Route("/users", func(r chi.Router) {
		r.Use(marshallingMiddleware)

		r.Post("/signup", h.signup)
	})
}

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h UserHandler) signup(w http.ResponseWriter, r *http.Request) {
	signupRequest := signupRequest{}
	err := bindRequest(r, &signupRequest)
	if err != nil {
		renderError(r.Context(), w, err)
		return
	}

	err = h.service.Signup(r.Context(), user.SignupInfo{
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	})
	if err != nil {
		renderError(r.Context(), w, errBusinessLogic(err))
		return
	}

	renderNoContent(r.Context(), w)
}
