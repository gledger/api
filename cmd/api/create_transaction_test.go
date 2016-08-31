package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func Test_decodeCreateTransactionRequest_ReturnsRequestData(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", "/",
		bytes.NewBufferString(
			`{"data":{"type":"transactions","attributes":{"payee":"Test Payee","amount":1,"occurred_at":"2016-08-25"},"relationships":{"account":{"data":{"id":"account-id"}},"envelope":{"data":{"id":"foo"}}}}}`,
		),
	)

	i, err := decodeCreateTransactionRequest(context.Background(), req)
	require.NoError(t, err)
	if r, ok := i.(createTransactionRequest); assert.True(t, ok) {
		assert.Equal(t, "Test Payee", r.Data.Attributes.Payee)
		assert.Equal(t, int64(1), r.Data.Attributes.Amount)
		assert.Equal(t, Date(time.Date(2016, 8, 25, 0, 0, 0, 0, time.UTC)), r.Data.Attributes.OccurredAt)
		assert.Equal(t, "account-id", r.Data.Relationships.Account.Data.ID)
		assert.Equal(t, "foo", r.Data.Relationships.Envelope.Data.ID)
	}
}

func Test_decodeCreateTransactionRequest_ErrorsWithMissingAttributes(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", "/",
		bytes.NewBufferString(
			`{}`,
		),
	)

	_, err := decodeCreateTransactionRequest(context.Background(), req)
	assert.Error(t, err)
}

func Test_decodeCreateTransactionRequest_ErrorsWithMissingData(t *testing.T) {
	for attr, e := range map[string]string{
		`{"payee":"","amount":1,"occurred_at":"2016-08-25"}`:               "Payee is required",
		`{"payee":"Test Payee","amount":null,"occurred_at":"2016-08-25"}`:  "Amount is required",
		`{"payee":"Test Payee","amount":"foo","occurred_at":"2016-08-25"}`: "json: cannot unmarshal string into Go value of type int64",
		`{"payee":"Test Payee","amount":1,"occurred_at":"foo"}`:            "parsing time \"foo\" as \"2006-01-02\": cannot parse \"foo\" as \"2006\"",
		`{"payee":"Test Payee","amount":1,"occurred_at":null}`:             "parsing time \"\" as \"2006-01-02\": cannot parse \"\" as \"2006\"",
	} {
		req, _ := http.NewRequest(
			"GET", "/",
			bytes.NewBufferString(
				fmt.Sprintf(`{"data":{"type":"transactions","attributes":%s},"relationships":{"account":{"data":{"id":"account-id"}},"envelope":{"data":{"id":"foo"}}}}`, attr),
			),
		)

		_, err := decodeCreateTransactionRequest(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, e, err.Error())
	}
}
