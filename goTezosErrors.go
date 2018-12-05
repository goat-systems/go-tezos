package goTezos


type NoClientError struct {

}

func (this NoClientError) Error() string {
	return "GoTezos did not find any healthy Tezos Node"
}