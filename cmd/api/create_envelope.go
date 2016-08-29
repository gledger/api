package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gledger/api"
	"github.com/go-kit/kit/endpoint"

	"golang.org/x/net/context"
)

func makeCreateEnvelopeEndpoint(svc gledger.EnvelopeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createEnvelopeRequest)
		e, err := svc.Create(gledger.Envelope{
			Name: req.Data.Attributes.Name,
		})

		return jsonAPIDocument{
			Data: jsonAPIEnvelopeResource{
				Type: "envelopes",
				ID:   e.UUID,
				Attributes: &jsonAPIEnvelopeResourceAttributes{
					Name:    e.Name,
					Balance: e.Balance,
				},
			},
		}, err
	}
}

type createEnvelopeRequest struct {
	Data jsonAPIEnvelopeResource `json:"data"`
}

func decodeCreateEnvelopeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createEnvelopeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.Data.Attributes == nil {
		return nil, errors.New("attributes are required")
	} else if request.Data.Attributes.Name == "" {
		return nil, errors.New("name attribute is required")
	}

	return request, nil
}
