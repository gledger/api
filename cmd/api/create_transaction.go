package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func makeCreateTransactionEndpoint(svc gledger.TransactionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTransactionRequest)
		t, err := svc.Create(gledger.Transaction{
			AccountUuid: req.AccountUuid,
			OccurredAt:  time.Time(req.OccurredAt),
			Payee:       req.Payee,
			Amount:      req.Amount,
			Cleared:     req.Cleared,
		})

		return jsonApiDocument{
			Data: jsonApiTransactionResource{
				Type: "transactions",
				Id:   t.Uuid,
				Attributes: jsonApiTransactionResourceAttributes{
					OccurredAt: Date(t.OccurredAt),
					Payee:      t.Payee,
					Amount:     t.Amount,
					Cleared:    t.Cleared,
					Reconciled: t.Reconciled,
				},
			},
		}, err
	}
}

type createTransactionRequest struct {
	AccountUuid string
	OccurredAt  Date   `json:"occurred_at"`
	Payee       string `json:"payee"`
	Amount      int64  `json:"amount"`
	Cleared     bool   `json:"cleared"`
}

func decodeCreateTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createTransactionRequest
	vars := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.AccountUuid = vars["uuid"]

	return request, nil
}
