package gorequests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RequestInstance struct {
	method string
	url    string
	json   interface{}
}

func Request() *RequestInstance {
	return new(RequestInstance)
}

func (r *RequestInstance) Method(method string) *RequestInstance {
	r.method = method
	return r
}

func (r *RequestInstance) Url(url string) *RequestInstance {
	r.url = url
	return r
}

func (r *RequestInstance) Json(json interface{}) *RequestInstance {
	r.json = json
	return r
}

func (r *RequestInstance) Exec() error {
	var bodyReader io.Reader
	if r.json != nil {
		body, err := json.Marshal(r.json)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = res
	return nil
}
