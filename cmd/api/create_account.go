package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func makeCreateAccountEndpoint(svc gledger.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		a, err := svc.Create(gledger.Account{
			Name:   req.Name,
			Type:   req.Type,
			Active: true,
		})

		return createAccountResponse{
			Uuid:   a.Uuid,
			Active: a.Active,
			createAccountRequest: createAccountRequest{
				Name: a.Name,
				Type: a.Type,
			},
		}, err
	}
}

type createAccountRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type createAccountResponse struct {
	Uuid string `json:"uuid"`

	createAccountRequest

	Active bool `json:"active"`
}
