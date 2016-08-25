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
		ts, err := svc.AllForAccount(req.AccountUUID)

		res := make([]jsonAPITransactionResource, len(ts), len(ts))
		for i, t := range ts {
			res[i] = jsonAPITransactionResource{
				ID:   t.UUID,
				Type: "transactions",
				Attributes: jsonAPITransactionResourceAttributes{
					Payee:        t.Payee,
					Amount:       t.Amount,
					RollingTotal: t.RollingTotal,
					OccurredAt:   Date(t.OccurredAt),
					Cleared:      t.Cleared,
					Reconciled:   t.Reconciled,
				},
				Relationships: jsonAPITransactionsRelationships{
					Account: jsonAPITransactionsRelationshipsAccount{
						Data: jsonAPIAccountResource{
							Type: "accounts",
							ID:   t.AccountUUID,
						},
					},
				},
			}
		}

		return jsonAPIDocument{Data: res}, err
	}
}

type readAccountTransactionsRequest struct {
	AccountUUID string
}

func decodeReadAccountTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return readAccountTransactionsRequest{
		AccountUUID: mux.Vars(r)["uuid"],
	}, nil
}
