package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func makeCreateAccountEndpoint(svc gledger.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		a, err := svc.Create(gledger.Account{
			Name:   req.Data.Attributes.Name,
			Type:   req.Data.Attributes.Type,
			Active: true,
		})

		return jsonAPIDocument{
			Data: jsonAPIAccountResource{
				Type: "accounts",
				ID:   a.UUID,
				Attributes: &jsonAPIAccountResourceAttributes{
					Name:   a.Name,
					Type:   a.Type,
					Active: a.Active,
				},
			},
		}, err
	}
}

type createAccountRequest struct {
	Data jsonAPIAccountResource `json:"data"`
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
