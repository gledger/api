package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func makeAllAccountsEndpoint(svc gledger.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		as, err := svc.All()

		res := make([]createAccountResponse, len(as), len(as))
		for i, a := range as {
			res[i] = createAccountResponse{
				Uuid:   a.Uuid,
				Active: a.Active,
				createAccountRequest: createAccountRequest{
					Name: a.Name,
					Type: a.Type,
				},
			}
		}

		return res, err
	}
}
