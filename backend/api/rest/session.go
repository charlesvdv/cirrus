package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/charlesvdv/cirrus/backend/pkg/identity"
)

type SessionService interface {
	Authenticate(ctx context.Context, credential identity.AuthenticationCredential) (identity.AuthenticationTokens, error)
}

func NewSessionHandler(service SessionService) SessionHandler {
	return SessionHandler{
		service: service,
	}
}

type SessionHandler struct {
	service SessionService
}

func (h SessionHandler) register(router *chi.Mux, ctx handlerContext) {
	router.Route("/session", func(r chi.Router) {
		r.Use(marshallingMiddleware)

		r.Post("/authenticate", h.authenticate)
	})
}

type authenticateRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ClientReference string `json:"client_reference"`
}

func responseToken(t identity.Token) token {
	return token{
		Token:     t.Token(),
		ExpiredAt: t.ExpiredAt(),
	}
}

type token struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type authenticateResponse struct {
	AccessToken     token  `json:"access_token"`
	RefreshToken    token  `json:"refresh_token"`
	ClientReference string `json:"client_reference"`
}

func (h SessionHandler) authenticate(w http.ResponseWriter, r *http.Request) {
	authenticateRequest := authenticateRequest{}
	err := bindRequest(r, &authenticateRequest)
	if err != nil {
		renderError(r.Context(), w, err)
	}

	tokens, err := h.service.Authenticate(r.Context(), identity.AuthenticationCredential{
		Email:           authenticateRequest.Email,
		Password:        authenticateRequest.Password,
		ClientReference: authenticateRequest.ClientReference,
	})
	if err != nil {
		renderError(r.Context(), w, convertIdentityErr(err))
		return
	}

	render(r.Context(), w, authenticateResponse{
		AccessToken:     responseToken(tokens.AccessToken),
		RefreshToken:    responseToken(tokens.RefreshToken),
		ClientReference: tokens.ClientReference,
	})
}
