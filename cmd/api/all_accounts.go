package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func makeAllAccountsEndpoint(svc gledger.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		as, err := svc.All()

		res := make([]jsonApiAccountResource, len(as), len(as))

		for i, a := range as {
			res[i] = jsonApiAccountResource{
				Type: "accounts",
				Id:   a.Uuid,
				Attributes: &jsonApiAccountResourceAttributes{
					Name:   a.Name,
					Type:   a.Type,
					Active: a.Active,
				},
			}
		}

		return jsonApiDocument{
			Data: res,
		}, err
	}
}
