package rest

import (
	"context"
	"net/http"

	"github.com/charlesvdv/cirrus/backend/pkg/identity"
	"github.com/go-chi/chi"
)

type UserService interface {
	Signup(ctx context.Context, info identity.SignupInfo) error
	GetUser(ctx context.Context, userID identity.UserID) (identity.User, error)
}

func NewUserHandler(service UserService) UserHandler {
	return UserHandler{
		service: service,
	}
}

type UserHandler struct {
	service UserService
}

func (h UserHandler) register(router *chi.Mux, ctx handlerContext) {
	router.Route("/users", func(r chi.Router) {
		r.Use(marshallingMiddleware)

		r.Post("/", h.createUser)
		r.With(ctx.authMiddleware).Get("/", h.getUserInfo)
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

	err = h.service.Signup(r.Context(), identity.SignupInfo{
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	})
	if err != nil {
		renderError(r.Context(), w, convertIdentityErr(err))
		return
	}

	renderNoContent(r.Context(), w)
}

type userInfoResponse struct {
	Email string `json:"email"`
}

func (h UserHandler) getUserInfo(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(r.Context(), getUserID(r.Context()))
	if err != nil {
		renderError(r.Context(), w, convertIdentityErr(err))
	}

	render(r.Context(), w, userInfoResponse{
		Email: user.Email(),
	})
}
