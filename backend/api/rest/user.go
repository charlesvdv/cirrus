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

		r.Post("/", h.createUser)
	})
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	createUserRequest := createUserRequest{}
	err := bindRequest(r, &createUserRequest)
	if err != nil {
		renderError(r.Context(), w, err)
		return
	}

	err = h.service.Signup(r.Context(), user.SignupInfo{
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	})
	if err != nil {
		renderError(r.Context(), w, errBusinessLogic(err))
		return
	}

	renderNoContent(r.Context(), w)
}
