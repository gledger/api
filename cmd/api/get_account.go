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
		a, err := svc.Read(req.AccountUuid)

		return jsonApiDocument{
			Data: jsonApiAccountResource{
				Type: "accounts",
				Id:   a.Uuid,
				Attributes: &jsonApiAccountResourceAttributes{
					Name:    a.Name,
					Type:    a.Type,
					Active:  a.Active,
					Balance: a.Balance,
				},
				Relationships: &jsonApiAccountResourceRelationships{
					Transactions: map[string]map[string]string{
						"links": map[string]string{
							"related": fmt.Sprintf("/accounts/%s/transactions", a.Uuid),
						},
					},
				},
			},
		}, err
	}
}

type readAccountRequest struct {
	AccountUuid string
}

func decodeReadAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return readAccountRequest{
		AccountUuid: mux.Vars(r)["uuid"],
	}, nil
}
