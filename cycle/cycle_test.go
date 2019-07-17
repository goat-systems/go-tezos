package cycle

import (
	"testing"

	"gotest.tools/assert"
)

func Test_GetCurrent(t *testing.T) {
	cases := []struct {
		want    int
		wantErr bool
	}{
		{
			want:    127,
			wantErr: false,
		},
	}

	for _, tc := range cases {

		cycleService := NewCycleService(&blockServiceMock{})

		cycle, err := cycleService.GetCurrent()
		if !tc.wantErr {
			assert.NilError(t, err)
			assert.Equal(t, tc.want, cycle)
		} else {
			assert.Assert(t, err != nil)
		}
	}
}
