package main

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/zombor/gledger"
	"golang.org/x/net/context"
)

func makeReadAccountTransactionsEndpoint(svc gledger.TransactionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(readAccountTransactionsRequest)
		ts, err := svc.AllForAccount(req.AccountUuid)

		res := make([]jsonApiTransactionResource, len(ts), len(ts))
		for i, t := range ts {
			res[i] = jsonApiTransactionResource{
				Id:   t.Uuid,
				Type: "transactions",
				Attributes: jsonApiTransactionResourceAttributes{
					Payee:      t.Payee,
					Amount:     t.Amount,
					OccurredAt: Date(t.OccurredAt),
					Cleared:    t.Cleared,
					Reconciled: t.Reconciled,
				},
				Relationships: jsonApiTransactionsRelationships{
					Account: jsonApiTransactionsRelationshipsAccount{
						Data: jsonApiAccountResource{
							Type: "accounts",
							Id:   t.AccountUuid,
						},
					},
				},
			}
		}

		return jsonApiDocument{Data: res}, err
	}
}

type readAccountTransactionsRequest struct {
	AccountUuid string
}

func decodeReadAccountTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return readAccountTransactionsRequest{
		AccountUuid: mux.Vars(r)["uuid"],
	}, nil
}
