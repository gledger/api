package main

import (
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := errors.Cause(err).Error()

	if e, ok := err.(httptransport.Error); ok {
		if nfe, ok := e.Err.(interface {
			NotFound() bool
		}); ok && nfe.NotFound() {
			code = http.StatusNotFound
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

type errorWrapper struct {
	Error string `json:"error"`
}
