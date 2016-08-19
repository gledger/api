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

		res := make([]createTransactionResponse, len(ts), len(ts))
		for i, t := range ts {
			res[i] = createTransactionResponse{
				Uuid:       t.Uuid,
				Reconciled: t.Reconciled,
				createTransactionRequest: createTransactionRequest{
					AccountUuid: t.AccountUuid,
					OccurredAt:  Date(t.OccurredAt),
					Payee:       t.Payee,
					Amount:      t.Amount,
					Cleared:     t.Cleared,
				},
			}
		}

		return res, err
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
