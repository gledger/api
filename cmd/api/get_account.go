package main

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/zombor/gledger"
	"golang.org/x/net/context"
)

func makeReadAccountEndpoint(svc gledger.AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(readAccountRequest)
		a, err := svc.Read(req.AccountUUID)

		return jsonAPIDocument{
			Data: jsonAPIAccountResource{
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
							"related": fmt.Sprintf("/accounts/%s/transactions", a.UUID),
						},
					},
				},
			},
		}, err
	}
}

type readAccountRequest struct {
	AccountUUID string
}

func decodeReadAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return readAccountRequest{
		AccountUUID: mux.Vars(r)["uuid"],
	}, nil
}
