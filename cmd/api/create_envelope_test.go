package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_decodeCreateEnvelopeRequest_ReturnsRequestData(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", "/",
		bytes.NewBufferString(
			`{"data":{"type":"envelopes","attributes":{"name":"Test Name","type":"income"}}}`,
		),
	)

	i, err := decodeCreateEnvelopeRequest(context.Background(), req)
	assert.NoError(t, err)
	if r, ok := i.(createEnvelopeRequest); assert.True(t, ok) {
		assert.Equal(t, "Test Name", r.Data.Attributes.Name)
	}
}

func Test_decodeCreateEnvelopeRequest_ErrorsWithMissingAttributes(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", "/",
		bytes.NewBufferString(
			`{}`,
		),
	)

	_, err := decodeCreateEnvelopeRequest(context.Background(), req)
	assert.Error(t, err)
}

func Test_decodeCreateEnvelopeRequest_ErrorsWithMissingData(t *testing.T) {
	for _, e := range []string{
		`{"name":"","type":"income"}`,
		`{"name":"test","type":""}`,
		`{"name":"test","type":"foo"}`,
	} {
		req, _ := http.NewRequest(
			"GET", "/",
			bytes.NewBufferString(
				fmt.Sprintf(`{"data":{"attributes":%s}}`, e),
			),
		)

		_, err := decodeCreateEnvelopeRequest(context.Background(), req)
		assert.Error(t, err)
	}
}
