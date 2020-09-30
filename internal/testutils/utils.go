package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// CheckErr -
func CheckErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), errContains)
		}
	} else {
		assert.Nil(t, err)
	}
}
