package gledger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EnvelopeCreateAssignsUuid_IfNotExists(t *testing.T) {
	e := Envelope{Name: "test"}
	var saveCalled bool
	saveErr := errors.New("boom")
	saveEnvelope := func(en Envelope) error {
		assert.Equal(t, e.Name, en.Name)
		assert.NotEmpty(t, en.UUID)
		saveCalled = true
		return saveErr
	}

	svc := envelopeService{save: saveEnvelope}
	_, err := svc.Create(e)
	assert.Equal(t, saveErr, err)
	assert.True(t, saveCalled)
}
