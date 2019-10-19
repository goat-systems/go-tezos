package network

import (
	"encoding/json"
	"testing"

	"gotest.tools/assert"

	tzc "github.com/DefinitelyNotAGoat/go-tezos/v2/client"
)

func Test_GetVersions(t *testing.T) {
	cases := []struct {
		tzclient tzc.TezosClient
		want     []byte
	}{
		{
			tzclient: &client{
				ReturnBody: goldenVersions,
			},
			want: goldenVersions,
		},
	}

	for _, tc := range cases {
		ns := NewNetworkService(tc.tzclient)
		versions, err := ns.GetVersions()
		assert.NilError(t, err)
		versionsWant := Versions{}
		versionsWant, err = versionsWant.unmarshalJSON(goldenVersions)
		assert.NilError(t, err)

		jsonHave, _ := json.Marshal(versions)
		jsonWant, _ := json.Marshal(versionsWant)

		assert.Equal(t, string(jsonHave), string(jsonWant))

	}
}

func Test_GetConstants(t *testing.T) {
	cases := []struct {
		tzclient tzc.TezosClient
		want     []byte
	}{
		{
			tzclient: &client{
				ReturnBody: goldenConstants,
			},
			want: goldenConstants,
		},
	}

	for _, tc := range cases {
		ns := NewNetworkService(tc.tzclient)
		constants, err := ns.GetConstants()
		assert.NilError(t, err)
		constantsWant := Constants{}
		constantsWant, err = constantsWant.unmarshalJSON(goldenConstants)
		assert.NilError(t, err)

		jsonHave, _ := json.Marshal(constants)
		jsonWant, _ := json.Marshal(constantsWant)

		assert.Equal(t, string(jsonHave), string(jsonWant))

	}
}

func Test_GetChainID(t *testing.T) {
	cases := []struct {
		tzclient tzc.TezosClient
		want     []byte
	}{
		{
			tzclient: &client{
				ReturnBody: goldenChainID,
			},
			want: goldenChainID,
		},
	}

	for _, tc := range cases {
		ns := NewNetworkService(tc.tzclient)
		chainID, err := ns.GetChainID()
		assert.NilError(t, err)
		chainIDWant, err := unmarshalString(goldenChainID)
		assert.NilError(t, err)

		jsonHave, _ := json.Marshal(chainID)
		jsonWant, _ := json.Marshal(chainIDWant)

		assert.Equal(t, string(jsonHave), string(jsonWant))

	}
}

func Test_Connections(t *testing.T) {
	cases := []struct {
		tzclient tzc.TezosClient
		want     []byte
	}{
		{
			tzclient: &client{
				ReturnBody: goldenConnections,
			},
			want: goldenConnections,
		},
	}

	for _, tc := range cases {
		ns := NewNetworkService(tc.tzclient)
		connections, err := ns.GetConnections()
		assert.NilError(t, err)

		var connectionsWant Connections
		connectionsWant, err = connectionsWant.unmarshalJSON(goldenConnections)
		assert.NilError(t, err)

		jsonHave, _ := json.Marshal(connections)
		jsonWant, _ := json.Marshal(connectionsWant)

		assert.Equal(t, string(jsonHave), string(jsonWant))

	}
}
