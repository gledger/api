package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/vnd.api+json")
	return json.NewEncoder(w).Encode(response)
}
