package gotezos

// func Test_B58cencode(t *testing.T) {
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
// 			res := B58cencode(tt.input.payload, tt.input.prefix)
// 			assert.Equal(t, tt.want.res, res)
// 		})
// 	}
// }
