package gotezos

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cycle(t *testing.T) {

	type input struct {
		handler http.Handler
		cycle   int
	}

	type want struct {
		err         bool
		errContains string
		cycle       Cycle
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"failed to get head block",
			input{
				gtGoldenHTTPMock(newBlockMock().handler([]byte(`not_block_data`), blankHandler)),
				10,
			},
			want{
				true,
				"could not get cycle '10': could not get head block",
				Cycle{},
			},
		},
		{
			"failed to get cycle because cycle is in the future",
			input{
				gtGoldenHTTPMock(newBlockMock().handler(mockBlockRandom, blankHandler)),
				300,
			},
			want{
				true,
				"request is in the future",
				Cycle{},
			},
		},
		{
			"failed to get block less than cycle",
			input{
				gtGoldenHTTPMock(
					newBlockMock().handler(
						mockBlockRandom,
						newBlockMock().handler(
							[]byte(`not_block_data`),
							blankHandler,
						),
					),
				),
				2,
			},
			want{
				true,
				"could not get block",
				Cycle{},
			},
		},
		{
			"failed to unmarshal cycle",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						[]byte(`bad_cycle_data`),
						newBlockMock().handler(
							mockBlockRandom,
							newBlockMock().handler(
								mockBlockRandom,
								blankHandler,
							),
						),
					),
				),
				2,
			},
			want{
				true,
				"could not unmarshal at cycle hash",
				Cycle{},
			},
		},
		{
			"failed to get cycle block level",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						mockCycle,
						newBlockMock().handler(
							mockBlockRandom,
							newBlockMock().handler(
								mockBlockRandom,
								newBlockMock().handler(
									[]byte(`not_block_data`),
									blankHandler,
								),
							),
						),
					),
				),
				2,
			},
			want{
				true,
				"could not get block",
				Cycle{
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
				},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(blankHandler)),
				2,
			},
			want{
				false,
				"",
				Cycle{
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
					BlockHash:    "BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			cycle, err := gt.Cycle(tt.input.cycle)
			if tt.want.err {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.errContains)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.cycle, cycle)
		})
	}
}

func mockCycleSuccessful(next http.Handler) http.Handler {
	var blockmock blockMock
	var oldHTTPBlock blockMock
	var blockAtLevel blockMock
	return cycleHandlerMock(
		mockCycle,
		blockmock.handler(
			mockBlockRandom,
			oldHTTPBlock.handler(
				mockBlockRandom,
				blockAtLevel.handler(
					mockBlockRandom,
					next,
				),
			),
		),
	)
}

func mockCycleFailed(next http.Handler) http.Handler {
	var blockmock blockMock
	var oldHTTPBlock blockMock
	return cycleHandlerMock(
		[]byte(`bad_cycle_data`),
		blockmock.handler(
			mockBlockRandom,
			oldHTTPBlock.handler(
				mockBlockRandom,
				blankHandler,
			),
		),
	)
}
