package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/gledger/api"
)

func makeAllEnvelopesEndpoint(svc gledger.EnvelopeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		es, err := svc.All()

		res := make([]jsonAPIEnvelopeResource, len(es), len(es))

		for i, e := range es {
			res[i] = jsonAPIEnvelopeResource{
				Type: "envelopes",
				ID:   e.UUID,
				Attributes: &jsonAPIEnvelopeResourceAttributes{
					Name:    e.Name,
					Type:    e.Type,
					Balance: e.Balance,
				},
			}
		}

		return jsonAPIDocument{
			Data: res,
		}, err
	}
}
