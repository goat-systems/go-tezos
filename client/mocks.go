package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type httpClientMock struct {
	ReturnStatus int
	ReturnBody   []byte
}

func (h *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader(h.ReturnBody)),
		StatusCode: h.ReturnStatus,
	}, nil
}

func (h *httpClientMock) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader(h.ReturnBody)),
		StatusCode: h.ReturnStatus,
	}, nil
}

func (h *httpClientMock) CloseIdleConnections() {

}

// func areEqualJSON(s1, s2 string) (bool, error) {
// 	var o1 interface{}
// 	var o2 interface{}

// 	var err error
// 	err = json.Unmarshal([]byte(s1), &o1)
// 	if err != nil {
// 		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
// 	}
// 	err = json.Unmarshal([]byte(s2), &o2)
// 	if err != nil {
// 		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
// 	}

// 	return reflect.DeepEqual(o1, o2), nil
// }
