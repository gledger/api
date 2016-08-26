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
			`{"data":{"type":"transactions","attributes":{"payee":"Test Payee","amount":1,"occurred_at":"2016-08-25"}}}`,
		),
	)

	i, err := decodeCreateTransactionRequest(context.Background(), req)
	require.NoError(t, err)
	if r, ok := i.(createTransactionRequest); assert.True(t, ok) {
		assert.Equal(t, "Test Payee", r.Data.Attributes.Payee)
		assert.Equal(t, int64(1), r.Data.Attributes.Amount)
		assert.Equal(t, Date(time.Date(2016, 8, 25, 0, 0, 0, 0, time.UTC)), r.Data.Attributes.OccurredAt)
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
	for _, e := range []string{
		`{"payee":"","amount":1,"occurred_at":"2016-08-25"}`,
		`{"payee":"Test Payee","amount":null,"occurred_at":"2016-08-25"}`,
		`{"payee":"Test Payee","amount":"foo","occurred_at":"2016-08-25"}`,
		`{"payee":"Test Payee","amount":1,"occurred_at":"foo"}`,
		`{"payee":"Test Payee","amount":1,"occurred_at":null}`,
	} {
		req, _ := http.NewRequest(
			"GET", "/",
			bytes.NewBufferString(
				fmt.Sprintf(`{"data":{"attributes":%s}}`, e),
			),
		)

		_, err := decodeCreateTransactionRequest(context.Background(), req)
		assert.Error(t, err)
	}
}
