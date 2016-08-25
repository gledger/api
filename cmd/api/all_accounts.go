package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/gledger/api"
)

func makeAllAccountsEndpoint(svc gledger.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		as, err := svc.All()

		res := make([]jsonAPIAccountResource, len(as), len(as))

		for i, a := range as {
			res[i] = jsonAPIAccountResource{
				Type: "accounts",
				ID:   a.UUID,
				Attributes: &jsonAPIAccountResourceAttributes{
					Name:    a.Name,
					Type:    a.Type,
					Active:  a.Active,
					Balance: a.Balance,
				},
			}
		}

		return jsonAPIDocument{
			Data: res,
		}, err
	}
}
