package gotezos

// import (
// 	"encoding/json"
// 	"testing"

// 	tzc "github.com/DefinitelyNotAGoat/go-tezos/v2/client"
// 	"gotest.tools/assert"
// )

// func Test_Get(t *testing.T) {
// 	cases := []struct {
// 		cycle    int
// 		tzclient tzc.TezosClient
// 		want     []byte
// 	}{
// 		{
// 			tzclient: &clientMock{
// 				ReturnBody: goldenGet,
// 			},
// 			want: goldenGet,
// 		},
// 	}

// 	for _, tc := range cases {
// 		sss := NewSnapshotService(
// 			&cycleServiceMock{},
// 			tc.tzclient,
// 			&blockServiceMock{},
// 			goldenConstants,
// 		)

// 		snapshot, err := sss.Get(tc.cycle)
// 		assert.NilError(t, err)

// 		jsonHave, _ := json.Marshal(snapshot)
// 		jsonWant, _ := json.Marshal(goldenSnapshot)

// 		assert.Equal(t, string(jsonHave), string(jsonWant))
// 	}
// }

// func Test_GetAll(t *testing.T) {
// 	cases := []struct {
// 		cycle    int
// 		tzclient tzc.TezosClient
// 		want     []byte
// 	}{
// 		{
// 			tzclient: &clientMock{
// 				ReturnBody: goldenGet,
// 			},
// 			want: goldenGet,
// 		},
// 	}

// 	for _, tc := range cases {
// 		sss := NewSnapshotService(
// 			&cycleServiceMock{},
// 			tc.tzclient,
// 			&blockServiceMock{},
// 			goldenConstants,
// 		)

// 		snapshot, err := sss.Get(tc.cycle)
// 		assert.NilError(t, err)

// 		jsonHave, _ := json.Marshal(snapshot)
// 		jsonWant, _ := json.Marshal(goldenSnapshot)

// 		assert.Equal(t, string(jsonHave), string(jsonWant))
// 	}
// }
