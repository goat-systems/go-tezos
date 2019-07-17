package snapshot

import (
	"github.com/DefinitelyNotAGoat/go-tezos/block"
	"github.com/DefinitelyNotAGoat/go-tezos/network"
)

var (
	goldenConstants = network.Constants{
		ProofOfWorkNonceSize:         8,
		NonceLength:                  32,
		MaxRevelationsPerBlock:       32,
		MaxOperationDataLength:       16384,
		MaxProposalsPerDelegate:      20,
		PreservedCycles:              3,
		BlocksPerCycle:               2048,
		BlocksPerCommitment:          32,
		BlocksPerRollSnapshot:        256,
		BlocksPerVotingPeriod:        8192,
		TimeBetweenBlocks:            []string{"30", "40"},
		EndorsersPerBlock:            32,
		HardGasLimitPerOperation:     "400000",
		HardGasLimitPerBlock:         "4000000",
		ProofOfWorkThreshold:         "70368744177663",
		TokensPerRoll:                "10000000000",
		MichelsonMaximumTypeSize:     1000,
		SeedNonceRevelationTip:       "125000",
		OriginationSize:              257,
		BlockSecurityDeposit:         "160000000",
		EndorsementSecurityDeposit:   "20000000",
		BlockReward:                  "16000000",
		EndorsementReward:            "2000000",
		CostPerByte:                  "1000",
		HardStorageLimitPerOperation: "60000",
	}

	goldenSnapshot = Snapshot{
		Cycle:           0,
		Number:          5,
		AssociatedHash:  "BMXVTnGN7rwaCE34yuAuKzTHaPgyCUBxuVkM2Bbfo5jZvrrbZrY",
		AssociatedBlock: 1,
	}

	goldenGet = []byte(`{"last_roll":[],"nonces":[],"random_seed":"6ac67b546fb98acacb8b5c435acff959217b85907c3e8762875ce3afc39dbab3","roll_snapshot":5}`)
)

type cycleServiceMock struct{}

func (c *cycleServiceMock) GetCurrent() (int, error) {
	return 9, nil
}

type clientMock struct {
	ReturnBody []byte
}

func (c *clientMock) Post(path, args string) ([]byte, error) {
	return c.ReturnBody, nil
}

func (c *clientMock) Get(path string, params map[string]string) ([]byte, error) {
	return c.ReturnBody, nil
}

type blockServiceMock struct {
}

func (b *blockServiceMock) GetHead() (block.Block, error) {
	return block.Block{
		Metadata: block.Metadata{
			Level: block.Level{
				Cycle: 9,
			},
		},
	}, nil
}

func (b *blockServiceMock) Get(id interface{}) (block.Block, error) {
	return block.Block{
		Hash: "BMXVTnGN7rwaCE34yuAuKzTHaPgyCUBxuVkM2Bbfo5jZvrrbZrY",
	}, nil
}

func (b *blockServiceMock) IDToString(id interface{}) (string, error) {
	return "", nil
}
