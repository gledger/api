package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/gledger/api"
)

func makeCreateTransactionEndpoint(svc gledger.TransactionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if req, ok := request.(createTransactionRequest); ok {
			t, err := svc.Create(gledger.Transaction{
				AccountUUID:  req.Data.Relationships.Account.Data.ID,
				EnvelopeUUID: req.Data.Relationships.Envelope.Data.ID,
				OccurredAt:   time.Time(req.Data.Attributes.OccurredAt),
				Payee:        req.Data.Attributes.Payee,
				Amount:       req.Data.Attributes.Amount,
				Cleared:      req.Data.Attributes.Cleared,
			})

			return jsonAPIDocument{
				Data: jsonAPITransactionResource{
					Type: "transactions",
					ID:   t.UUID,
					Attributes: jsonAPITransactionResourceAttributes{
						OccurredAt: Date(t.OccurredAt),
						Payee:      t.Payee,
						Amount:     t.Amount,
						Cleared:    t.Cleared,
						Reconciled: t.Reconciled,
					},
					Relationships: jsonAPITransactionsRelationships{
						Account: jsonAPITransactionsRelationshipsAccount{
							Data: jsonAPIAccountResource{
								Type: "accounts",
								ID:   t.AccountUUID,
							},
						},
						Envelope: jsonAPITransactionsRelationshipsEnvelope{
							Data: jsonAPIEnvelopeResource{
								Type: "Envelopes",
								ID:   t.EnvelopeUUID,
							},
						},
					},
				},
			}, err
		}

		return nil, errors.New("request was not createTransactionRequest")
	}
}

func decodeCreateTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createTransactionRequest

	vars := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.Data.Type = "transactions"
	request.Data.Relationships.Account.Data.ID = vars["uuid"]
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
	if request.Data.Relationships.Envelope.Data.ID == "" {
		return request, errors.New("`envelope` relationship `id` is required")
	}

	return request, nil
}

type createTransactionRequest struct {
	Data jsonAPITransactionResource `json:"data"`
}
