package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_secp256k1Curve_getPublicKey(t *testing.T) {
	curve := &secp256k1Curve{}

	type want struct {
		wantErr bool
		pubKey  []byte
	}

	type input struct {
		privateKey []byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				privateKey: []byte{162, 200, 103, 46, 138, 63, 231, 67, 94, 53, 150, 55, 123, 94, 78, 228, 227, 233, 72, 38, 31, 241, 95, 29, 235, 93, 26, 31, 30, 196, 220, 216},
			},
			want{
				wantErr: false,
				pubKey:  []byte{3, 42, 15, 71, 222, 26, 129, 217, 176, 186, 177, 3, 70, 82, 109, 165, 209, 229, 84, 16, 19, 204, 162, 50, 54, 112, 26, 46, 177, 5, 182, 129, 255},
			},
		},
		{
			"is successful with len(Y) < 32",
			input{
				privateKey: []byte{214, 48, 92, 109, 221, 55, 16, 27, 97, 225, 74, 13, 58, 195, 209, 210, 104, 89, 190, 164, 218, 10, 252, 244, 194, 205, 248, 176, 147, 7, 128, 245},
			},
			want{
				wantErr: false,
				// len(pubKey.Y.Bytes()) == 30
				pubKey: []byte{3, 98, 55, 112, 26, 254, 221, 136, 82, 122, 229, 227, 42, 115, 73, 116, 184, 103, 48, 130, 181, 65, 165, 59, 153, 186, 194, 175, 153, 207, 206, 92, 18},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			pubKey, err := curve.getPublicKey(tt.input.privateKey)
			if tt.want.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.pubKey, pubKey)
		})
	}
}
