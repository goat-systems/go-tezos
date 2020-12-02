package keys

import (
	"testing"

	"github.com/goat-systems/go-tezos/v4/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_FromEncryptedSecret(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		secretKey   string
		publicKey   string
		address     string
	}

	type input struct {
		Esk      string
		Password string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful with edesk",
			input{
				Esk:      "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
				Password: "password12345##",
			},
			want{
				false,
				"",
				"edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF",
				"edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5",
				"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
			},
		},
		{
			"is successful with spesk",
			input{
				Esk:      "spesk24LjVwuCRhGsFYPGATnwaHAw7eZ6phDvGPntSHrSEVYstNpC3Zuq5k7oHTE1pAkVifNtJ1XW5UwCYcJC5BZ",
				Password: "abcd1234",
			},
			want{
				false,
				"",
				"spsk2psNeAQ88pKFnZikoZNb37zRbDmaGgQUtYrwwJZT3RcUspwL7N",
				"sppk7ZZADMMS4cwsu3odb7BAu9mx3DZYHmXWWL9GNhKremaJXqytGBc",
				"tz2TUwYWy5VP7ChX2xjXtGxxdfCnEQsotdeQ",
			},
		},
		{
			"is successful with p2esk",
			input{
				Esk:      "p2esk2UKLR5vLjQFjrotvStrYS6SXUFmGt9fkRjFydFwvFBJcbxWESvjfXtYKesgZaBHA7dx9MZJSwknhakFPZoc",
				Password: "abcd1234",
			},
			want{
				false,
				"",
				"p2sk3UumbKMrb6Wo1Jm5qTSMhUrCyAFTK4LMWgVma9njNLGc2Wcx9S",
				"p2pk6594Hd4VEVPydvK67c2GVikNXWjLiv2tkPUVvd8XMAXqd4CYxdK",
				"tz3fU9apdFnzoPhi4LB8AdxoiSVwLYM4kQ1F",
			},
		},
		// {
		// 	"is successful with mnemonic",
		// 	NewKeyInput{
		// 		Kind:     Ed25519,
		// 		Mnemonic: "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
		// 		Password: "PYh8nXDQLB",
		// 		Email:    "vksbjweo.qsrgfvbw@tezos.example.org",
		// 	},
		// 	want{
		// 		false,
		// 		"",
		// 		"edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
		// 		"edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
		// 		"tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
		// 	},
		// },
		// {
		// 	"is successful with base58",
		// 	NewKeyInput{
		// 		Kind:          Ed25519,
		// 		EncodedString: "edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
		// 	},
		// 	want{
		// 		false,
		// 		"",
		// 		"edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
		// 		"edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
		// 		"tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
		// 	},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			key, err := FromEncryptedSecret(tt.input.Esk, tt.input.Password)
			testutils.CheckErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.secretKey, key.GetSecretKey())
			assert.Equal(t, tt.want.publicKey, key.PubKey.GetPublicKey())
			assert.Equal(t, tt.want.address, key.PubKey.GetAddress())
		})
	}
}

func Test_FromBytes(t *testing.T) {
	privKey := []byte{117, 121, 196, 136, 31, 185, 152, 208, 67, 65, 123, 124, 4, 88, 42, 161, 81, 121, 241, 37, 197, 48, 62, 30, 229, 106, 150, 120, 3, 77, 149, 176}
	key, err := FromBytes(privKey, Ed25519)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF", key.GetSecretKey())
	assert.Equal(t, "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5", key.PubKey.GetPublicKey())
	assert.Equal(t, "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo", key.PubKey.GetAddress())
}

func Test_FromHex(t *testing.T) {
	privKey := "7579c4881fb998d043417b7c04582aa15179f125c5303e1ee56a9678034d95b0"
	key, err := FromHex(privKey, Ed25519)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF", key.GetSecretKey())
	assert.Equal(t, "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5", key.PubKey.GetPublicKey())
	assert.Equal(t, "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo", key.PubKey.GetAddress())
}

// func Test_FromBase64(t *testing.T) {
// 	privKey := base64.EncodeToString([]byte{117, 121, 196, 136, 31, 185, 152, 208, 67, 65, 123, 124, 4, 88, 42, 161, 81, 121, 241, 37, 197, 48, 62, 30, 229, 106, 150, 120, 3, 77, 149, 176})
// 	key, err := FromBase64(privKey, Ed25519)
// 	testutils.CheckErr(t, false, "tt.want.containsErr", err)
// 	assert.Equal(t, "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF", key.GetSecretKey())
// 	assert.Equal(t, "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5", key.PubKey.GetPublicKey())
// 	assert.Equal(t, "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo", key.PubKey.GetPublicKeyHash())
// }

func Test_FromBase58(t *testing.T) {
	privKey := "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF"
	key, err := FromBase58(privKey, Ed25519)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF", key.GetSecretKey())
	assert.Equal(t, "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5", key.PubKey.GetPublicKey())
	assert.Equal(t, "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo", key.PubKey.GetAddress())
}

func Test_FromMnemonic(t *testing.T) {
	type want struct {
		wantErr     bool
		containsErr string
		secretKey   string
		publicKey   string
		address     string
	}

	type input struct {
		mnemonic string
		email    string
		password string
		kind     ECKind
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful with mnemonic",
			input{
				kind:     Ed25519,
				mnemonic: "normal dash crumble neutral reflect parrot know stairs culture fault check whale flock dog scout",
				password: "PYh8nXDQLB",
				email:    "vksbjweo.qsrgfvbw@tezos.example.org",
			},
			want{
				false,
				"",
				"edskRxB2DmoyZSyvhsqaJmw5CK6zYT7dbkUfEVSiQeWU1gw3ZMnC99QMMXru3imsbUrLhvuHktrymvNqhMxkhz7Y4LJAtevW5V",
				"edpkvEoAbkdaGALxi2FfeefB8hUkMZ4J1UVwkzyumx2GvbVpkYUHnm",
				"tz1Qny7jVMGiwRrP9FikRK95jTNbJcffTpx1",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			key, err := FromMnemonic(tt.input.mnemonic, tt.input.email, tt.input.password, tt.input.kind)
			testutils.CheckErr(t, tt.want.wantErr, tt.want.containsErr, err)
			assert.Equal(t, tt.want.secretKey, key.GetSecretKey())
			assert.Equal(t, tt.want.publicKey, key.PubKey.GetPublicKey())
			assert.Equal(t, tt.want.address, key.PubKey.GetAddress())
		})
	}
}
