package gledger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateAssignsUuid_IfNotExists(t *testing.T) {
	a := Account{Name: "test"}
	var saveCalled bool
	saveErr := errors.New("boom")
	saveAccount := func(act Account) error {
		assert.Equal(t, a.Name, act.Name)
		assert.NotEmpty(t, act.UUID)
		saveCalled = true
		return saveErr
	}

	svc := accountService{saveAccount: saveAccount}
	_, err := svc.Create(a)
	assert.Equal(t, saveErr, err)
	assert.True(t, saveCalled)
}
