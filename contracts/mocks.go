package contracts

var (
	goldenStorage = []byte(`{Elt "tz1gH29qAVaNfv7imhPthCwpUBcqmMdLWxPG" (Pair "Jackson" (Pair 100000 23))}`)
)

type clientMock struct {
	ReturnBody []byte
}

func (c *clientMock) Post(path, args string) ([]byte, error) {
	return c.ReturnBody, nil
}

func (c *clientMock) Get(path string, params map[string]string) ([]byte, error) {
	return c.ReturnBody, nil
}
