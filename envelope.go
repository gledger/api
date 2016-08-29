package gledger

import "github.com/twinj/uuid"

type Envelope struct {
	UUID, Name string
	Balance    int64
}

type EnvelopeService interface {
	Create(Envelope) (Envelope, error)
	All() ([]Envelope, error)
}

func NewEnvelopeService(
	s func(Envelope) error,
	a func() ([]Envelope, error),
) envelopeService {
	return envelopeService{save: s, all: a}
}

type envelopeService struct {
	save func(Envelope) error
	all  func() ([]Envelope, error)
}

func (es envelopeService) Create(e Envelope) (Envelope, error) {
	if e.UUID == "" {
		e.UUID = uuid.NewV4().String()
	}

	return e, es.save(e)
}

func (es envelopeService) All() ([]Envelope, error) {
	return es.all()
}
