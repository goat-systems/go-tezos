package gotezos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// func Test_b58cencode(t *testing.T) {
// 	type input struct {
// 		payload []byte
// 		prefix  prefix
// 	}

// 	type want struct {
// 		res string
// 	}

// 	cases := []struct {
// 		name  string
// 		input input
// 		want  want
// 	}{
// 		{
// 			"is successful with tz1",
// 			input{
// 				[]byte{117, 121, 196, 136, 31, 185, 152, 208, 67, 65, 123, 124, 4, 88, 42, 161, 81, 121, 241, 37, 197, 48, 62, 30, 229, 106, 150, 120, 3, 77, 149, 176, 84, 76, 85, 33, 188, 5, 113, 64, 14, 24, 19, 168, 43, 33, 121, 69, 55, 148, 148, 61, 195, 162, 152, 248, 170, 81, 226, 154, 199, 64,76, 163}
// 			},
// 			want{},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			res := b58cencode(tt.input.payload, tt.input.prefix)
// 			assert.Equal(t, tt.want.res, res)
// 		})
// 	}
// }

func TestCheckAddrFormat(t *testing.T) {
	correctAddr := "tz1buwfQ3j7gTSM5QU8bmG2YnfH8zEnsjm92"
	assert.True(t, CheckAddr(correctAddr))
	faultAddr := "tz1buwfQ3j7gTSM5QU8bmG2YnfH8zEnqjm92"
	assert.False(t, CheckAddr(faultAddr))

	correctAddr = "KT1TPBnfPq7XpCjL11HTzAeBkyxrSxUuskS9"
	assert.True(t, CheckAddr(correctAddr))
	faultAddr = "KT1TPBnfPq7XpCjL11HTzAeBkyxrSxUuskS8"
	assert.False(t, CheckAddr(faultAddr))
}

func TestCheckPubKey(t *testing.T) {
	correctPubKey := "edpkuKSXhRG9t6F3iSghmf5dMw1X6CKv9SHfGhEaPEf4miz7VDcLmy"
	assert.True(t, CheckPubKey(correctPubKey))
	faultPubKey := "edpkuKSXhRG9t6F3iSghmf5dMw1X6CKv9SHfGhEaPEf3miz7VDcLmy"
	assert.False(t, CheckPubKey(faultPubKey))
}

func TestEncodeSignature(t *testing.T) {
	hexSig := []string{
		"9039dc3ac3533f73081f3e3b66974f31a90175eedcced27bff86fc2258c022bca062a19246a43036d48f7a742de5a31d6843074bffdfd9cbb364a0da45f04a07",
		"25a4173f0077795c9a9f9da114e3fd9bf26969009c8af2fe65bde5077e66a717f891ca84c7c34514d3ee2cb5a7c66fc11d910f1d4d76a2f0703d8e72be9ae603",
		"e4fb50aa20227ee6df9c9efd7faf84c1a8ff455258021e822dc4498aa0f929ebc85d9b3259f539f4455bd8bcea789f39fe1b2e13c59ea7f500530296b6bac706",
	}
	expected := []string{
		"edsigtrgA4A8ydSRTh47MgH7CKuZtYEUAfo1c4gp2S4BeurVk1MQq7TTJe8UBQQesSZT5M6PBt2VfnkWrW2UCRmFMjmCkZaNeum",
		"edsigtcjPEgLjhSFDAD9pzVivuPcDHAHCAzXny5A4L5XCPkURJ2xgHuywfy4NVZrYcikw7YE1XVRWGuzADP9iibETdJMmqpTaGk",
		"edsigu3mHd5LcQawkxG4TzSDjxit2QAJudgutWqtujSb8PfXmN76LSnEZYrvP1vGTV1ZRPVvBgY8D8ra7xwaY8F8skK7Rtetdi5",
	}
	for index, sig := range hexSig {
		signature, err := EncodeSignature(sig)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected[index], signature)
	}
}
