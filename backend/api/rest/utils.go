package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func getMarshaller(ctx context.Context) marshaller {
	return ctx.Value(ctxKeyMarshaller).(marshaller)
}

func bindRequest(r *http.Request, v interface{}) error {
	marshaller := getMarshaller(r.Context())

	err := marshaller.unmarshal(r.Body, v)
	if err != nil {
		return errBadRequestErr(err, "Invalid body")
	}

	return nil
}

func renderError(ctx context.Context, w http.ResponseWriter, baseErr error) {
	marshaller := getMarshaller(ctx)

	var restError *restError
	if !errors.As(baseErr, &restError) {
		log.Ctx(ctx).Warn().Err(baseErr).Msg("Unknown error")
		renderError(ctx, w, errInternalError("Unknown error"))
		return
	}

	w.WriteHeader(restError.status)

	err := marshaller.marshal(w, restError)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msgf("Failed to marshal error: %#v. Sending an HTTP 500 with no body instead", restError)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Ctx(ctx).Debug().Err(restError).Msg("Returning error to user")
}

func render(ctx context.Context, w http.ResponseWriter, v interface{}) {
	renderWithStatus(ctx, w, http.StatusOK, v)
}

func renderWithStatus(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	marshaller := getMarshaller(ctx)

	err := marshaller.marshal(w, v)
	if err != nil {
		renderError(ctx, w, errInternalErrorErr(fmt.Errorf("failed to marshal response: %w", err), ""))
		return
	}

	w.WriteHeader(status)
}

func renderNoContent(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
