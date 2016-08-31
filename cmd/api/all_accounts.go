package main

import (
	"fmt"

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
				Relationships: &jsonAPIAccountResourceRelationships{
					Transactions: map[string]map[string]string{
						"links": map[string]string{
							"self": fmt.Sprintf("/accounts/%s/transactions", a.UUID),
						},
					},
				},
			}
		}

		return jsonAPIDocument{
			Data: res,
		}, err
	}
}
