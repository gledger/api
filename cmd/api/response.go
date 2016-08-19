package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
