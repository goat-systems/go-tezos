package keys

import (
	"encoding/hex"
	"testing"

	tzcrypt "github.com/completium/go-tezos/v4/internal/crypto"
	"github.com/completium/go-tezos/v4/internal/testutils"
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

func Test_PubKey(t *testing.T) {
	v, err := hex.DecodeString("0000861299624c9a3b52be10762c64bac282b1c02316")
	testutils.CheckErr(t, false, "", err)

	address := tzcrypt.B58cencode(v[2:], []byte{6, 161, 159})
	assert.Equal(t, "tz1XrwX7i9Nzh8e6UmG3VnFkAeoyWdTqDf3U", address)
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
	key, err := FromBase58(privKey)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "edskRsPBsKuULoLTEQV2R9UbvSZbzFqvoESvp1mYyQJU8xi9mJamt88r5uTXbWQpVHjSiPWWtnoyqTCuSLQLxbEKUXfwwTccsF", key.GetSecretKey())
	assert.Equal(t, "edpkuHMDkMz46HdRXYwom3xRwqk3zQ5ihWX4j8dwo2R2h8o4gPcbN5", key.PubKey.GetPublicKey())
	assert.Equal(t, "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo", key.PubKey.GetAddress())
}

func Test_FromBase58Pk(t *testing.T) {
	pubKey := "edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrpvKXhjn"
	pk, err := FromBase58Pk(pubKey)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, "edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrpvKXhjn", pk.GetPublicKey())
	assert.Equal(t, "tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb", pk.GetAddress())
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

func Test_Sign(t *testing.T) {
	k, _ := FromBase58("edsk3QoqBuvdamxouPhin7swCvkQNgq4jP5KZPbwWNnwdZpSpJiEbq")
	sigHex, _ := k.SignDataHex("050000")
	sigBytes, _ := k.SignDataBytes([]byte{5, 0, 0})
	assert.Equal(t, "edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W8vG", sigHex.ToBase58())
	assert.Equal(t, "edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W8vG", sigBytes.ToBase58())
}

func Test_Pkh(t *testing.T) {
	inputtz1, _ := hex.DecodeString("00009ecba6b6222ce85e64290ba325271857ded2e646")
	pkhtz1, _ := GetPkhFromBytes(inputtz1)
	assert.Equal(t, "tz1a7fRN1NZt84fuwGMpgrS4WYYiAbQEEozy", pkhtz1)

	inputKT1, _ := hex.DecodeString("013c783e38e2abe1c1d8f49d1851a4151f0573916b00")
	pkhKT1, _ := GetPkhFromBytes(inputKT1)
	assert.Equal(t, "KT1E6W9ugnHT1rsmTz8sERYpeSMFJnzQmPL6", pkhKT1)
}

func Test_isValidPkh(t *testing.T) {
	assert.Equal(t, true, IsValidPkh("tz1a7fRN1NZt84fuwGMpgrS4WYYiAbQEEozy"))
	assert.Equal(t, true, IsValidPkh("tz2BFTyPeYRzxd5aiBchbXN3WCZhx7BqbMBq"))
	assert.Equal(t, true, IsValidPkh("tz3hFR7NZtjT2QtzgMQnWb4xMuD6yt2YzXUt"))
	assert.Equal(t, true, IsValidPkh("KT1RXgG7X3wvx1GWv7aGdaALwbMsuDj8cXQe"))
	assert.Equal(t, false, IsValidPkh("blabla"))
	assert.Equal(t, false, IsValidPkh("tz1a7fRN1NZt84fuwGMpgrS4WYYiAbQE"))     // bad size
	assert.Equal(t, false, IsValidPkh("tz1a7fRN1NZt84fuwGMpgrS4WYYiAbQEEozl")) // invalid char
}

func Test_isValidPk(t *testing.T) {
	assert.Equal(t, true, IsValidPk("edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrpvKXhjn"))
	assert.Equal(t, true, IsValidPk("sppk7bcmsCiZmrzrfGpPHnZMx73s6pUC4Tf1zdASQ3rgXfq8uGP3wgV"))
	assert.Equal(t, true, IsValidPk("p2pk66tTYL5EvahKAXncbtbRPBkAnxo3CszzUho5wPCgWauBMyvybuB"))
	assert.Equal(t, false, IsValidPk("edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrpvKXhj"))
	assert.Equal(t, false, IsValidPk("edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrp"))       // bad size
	assert.Equal(t, false, IsValidPk("edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrpvKXhjl")) // invalid char
}

func Test_isValidSignature(t *testing.T) {
	assert.Equal(t, true, IsValidSignature("edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W8vG"))
	assert.Equal(t, true, IsValidSignature("spsig1VrEwwc2UC4v9v3oYJ96VwiKwdVKK7ZYdMs4JVWNtfj11sRz9RkvPBtCHMiG1LEp44PJBXDh7bAzpDjGoX4bH7heoPuGqa"))
	assert.Equal(t, true, IsValidSignature("p2sigZehGEs7pMMZCYhxDzRBbZkyWhBX26ctJ4BPCwGEV1CnEDpVq6DjbcUAxThDj6KKoMxpwTqvvaKs38pJb2mnb5rB8U3G9o"))
	assert.Equal(t, false, IsValidSignature("edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W8v"))
	assert.Equal(t, false, IsValidSignature("edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W"))    // bad size
	assert.Equal(t, false, IsValidSignature("edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W8vl")) // invalid char
}

func Test_isValidBlockHash(t *testing.T) {
	assert.Equal(t, true, IsValidBlockHash("BLLqDP5o6xTWzi9WEUJSxXw9qwW8qHd9jTCA9GeepPRmbV71Xxm"))
	assert.Equal(t, false, IsValidBlockHash("BLLqDP5o6xTWzi9WEUJSxXw9qwW8qHd9jTCA9GeepPRmbV71Xx"))  // bad size
	assert.Equal(t, false, IsValidBlockHash("BLLqDP5o6xTWzi9WEUJSxXw9qwW8qHd9jTCA9GeepPRmbV71Xx0")) // invalid char
}

func Test_isValidOperationHash(t *testing.T) {
	assert.Equal(t, true, IsValidOperationHash("op32G4toHNqGwDdEV9ZqwnrfdF244rjyZRfxrmx3zFxHwQCkwxv"))
	assert.Equal(t, true, IsValidOperationHash("oo5cw72WCRvfF2kCYPDYBGtPq4MjjprWruD9M2BLHuvGtpPcc8v"))
	assert.Equal(t, true, IsValidOperationHash("ootpLxzp6bQ8qXzayY6tMKEzJoMbC8Y4GHZW4PaG6FvyPTBRmnz"))
	assert.Equal(t, true, IsValidOperationHash("onws3e5BJPmCo8jGoLNFMmeLLmhkBjYdBPQUmsXweof3qtKmzGg"))
	assert.Equal(t, false, IsValidOperationHash("onws3e5BJPmCo8jGoLNFMmeLLmhkBjYdBPQUmsXweof3qtKmzG"))  // bad size
	assert.Equal(t, false, IsValidOperationHash("onws3e5BJPmCo8jGoLNFMmeLLmhkBjYdBPQUmsXweof3qtKmzG0")) // invalid char
}

func Test_CheckSignatureTz1(t *testing.T) {
	data := "050000"

	pk, _ := FromBase58Pk("edpkvGfYw3LyB1UcCahKQk4rF2tvbMUk8GFiTuMjL75uGXrpvKXhjn")

	sig := "edsigthXYBNW7i5E1WNd87fBRJKacJjK5amJVKcyXd6fGxmnQo2ESmmdgN6qJXgbUVJDXha8xi96r9GqjsPorWWpPEwXNG3W8vG"
	res, err := pk.CheckSignature(data, sig)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, true, res)

	// Same user, different data
	sigErr0 := "edsigu4S83X6eLC2e7pWsAHeGsd8bqFM3ADuroq7kHDCFLYG53hsWLJjxzaoT3uHhD2z29hdQmJHWsWnfhhXSH46msWtdzFHd3T"
	resErr0, err := pk.CheckSignature(data, sigErr0)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, false, resErr0)

	// Different user, same data
	sigErr1 := "edsigtqAGRMDM8hR4aGewGKfMei6eHkqUFMeCg7qitcyrQTCCYLVn5AJnPCj5JoFL4zzQmw6BnM25UmrpWvk9V31cUHWcS13ba2"
	resErr1, err := pk.CheckSignature(data, sigErr1)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, false, resErr1)
}

func Test_CheckSignatureTz2(t *testing.T) {
	data := "050000"

	pk, _ := FromBase58Pk("sppk7b4TURq2T9rhPLFaSz6mkBCzKzfiBjctQSMorvLD5GSgCduvKuf")

	sig := "spsig1VrEwwc2UC4v9v3oYJ96VwiKwdVKK7ZYdMs4JVWNtfj11sRz9RkvPBtCHMiG1LEp44PJBXDh7bAzpDjGoX4bH7heoPuGqa"
	res, err := pk.CheckSignature(data, sig)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, true, res)

	// Same user, different data
	sigErr0 := "spsig179FXkSPGX3ng3AH2tdDP4u8TcgXnHY2b1tj5vVmWFdRcG5KmgqvSCynxYJ7Gs8BBM2NW5z6raq7Up4hkSpBjQ2cPrVBzy"
	resErr0, err := pk.CheckSignature(data, sigErr0)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, false, resErr0)
}

func Test_CheckSignatureTz3(t *testing.T) {
	data := "050000"

	pk, _ := FromBase58Pk("p2pk65zwHGP9MdvANKkp267F4VzoKqL8DMNpPfTHUNKbm8S9DUqqdpw")

	sig := "p2sigZehGEs7pMMZCYhxDzRBbZkyWhBX26ctJ4BPCwGEV1CnEDpVq6DjbcUAxThDj6KKoMxpwTqvvaKs38pJb2mnb5rB8U3G9o"
	res, err := pk.CheckSignature(data, sig)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, true, res)

	// Same user, different data
	sigErr0 := "p2sigt6h7tWxS5Pz7D1qgfkRVemy4XuYGWw9Vo94kv7R6YVxfdtHSWeZ8QvkUfJbnmzqmK13fbBVQ89m7Tu11FpwDtEdH7jckY"
	resErr0, err := pk.CheckSignature(data, sigErr0)
	testutils.CheckErr(t, false, "", err)
	assert.Equal(t, false, resErr0)
}
