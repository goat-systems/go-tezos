package ledger

import "fmt"

// https://github.com/obsidiansystems/ledger-app-tezos/blob/master/APDUs.md
var (
	InsVersion             = []byte{0x00}
	InsGetPublicKey        = []byte{0x02}
	InsPromptPublicKey     = []byte{0x03}
	InsSign                = []byte{0x04}
	InsSignUnsafe          = []byte{0x05}
	CLA                    = []byte{0x80}
	INS_GET_PUBLIC_KEY     = []byte{0x02}
	INS_PROMPT_PUBLIC_KEY  = []byte{0x03}
	INS_SIGN               = []byte{0x04}
	FIRST_MESSAGE_SEQUENCE = []byte{0X00}
	LAST_MESSAGE_SEQUENCE  = []byte{0X81}
	OTHER_MESSAGE_SEQUENCE = []byte{0X01}
)

func HDPathTemplate(account int) string {
	return fmt.Sprintf("44'/1729'/%d'/0'", account)
}

// export const HDPathTemplate = (account: number) => {
// 	return `44'/1729'/${account}'/0'`;
//   };
