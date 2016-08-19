package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func makeCreateTransactionEndpoint(svc gledger.TransactionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if req, ok := request.(createTransactionRequest); ok {
			t, err := svc.Create(gledger.Transaction{
				AccountUuid: req.Data.Relationships.Account.Data.Id,
				OccurredAt:  time.Time(req.Data.Attributes.OccurredAt),
				Payee:       req.Data.Attributes.Payee,
				Amount:      req.Data.Attributes.Amount,
				Cleared:     req.Data.Attributes.Cleared,
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
					Relationships: jsonApiTransactionsRelationships{
						Account: jsonApiTransactionsRelationshipsAccount{
							Data: jsonApiAccountResource{
								Type: "accounts",
								Id:   t.AccountUuid,
							},
						},
					},
				},
			}, err
		} else {
			return nil, errors.New("request was not createTransactionRequest")
		}
	}
}

func decodeCreateTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createTransactionRequest

	vars := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.Data.Type = "transactions"
	request.Data.Relationships.Account.Data.Id = vars["uuid"]
	request.Data.Relationships.Account.Data.Type = "accounts"

	if request.Data.Attributes.Payee == "" {
		return request, errors.New("Payee is required")
	}
	if request.Data.Attributes.Amount == 0 {
		return request, errors.New("Amount is required")
	}
	if request.Data.Attributes.OccurredAt == Date(time.Time{}) {
		return request, errors.New("OccurredAt is required")
	}

	return request, nil
}

type createTransactionRequest struct {
	Data jsonApiTransactionResource `json:"data"`
}
