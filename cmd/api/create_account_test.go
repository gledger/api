package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_decodeCreateAccountRequest_ReturnsRequestData(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", "/",
		bytes.NewBufferString(
			`{"data":{"type":"accounts","attributes":{"name":"Test Name","type":"Test Type","active":true}}}`,
		),
	)

	i, err := decodeCreateAccountRequest(context.Background(), req)
	assert.NoError(t, err)
	if r, ok := i.(createAccountRequest); assert.True(t, ok) {
		assert.Equal(t, "Test Name", r.Data.Attributes.Name)
		assert.Equal(t, "Test Type", r.Data.Attributes.Type)
		assert.Equal(t, true, r.Data.Attributes.Active)
	}
}

func Test_decodeCreateAccountRequest_ErrorsWithMissingAttributes(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", "/",
		bytes.NewBufferString(
			`{}`,
		),
	)

	_, err := decodeCreateAccountRequest(context.Background(), req)
	assert.Error(t, err)
}

func Test_decodeCreateAccountRequest_ErrorsWithMissingData(t *testing.T) {
	for _, e := range []string{
		`{"name":"","type":"Test Type","active":true}`,
		`{"name":"Test Name","type":"","active":true}`,
	} {
		req, _ := http.NewRequest(
			"GET", "/",
			bytes.NewBufferString(
				fmt.Sprintf(`{"data":{"attributes":%s}}`, e),
			),
		)

		_, err := decodeCreateAccountRequest(context.Background(), req)
		assert.Error(t, err)
	}
}
