package gorequests

import (
	"encoding/json"
	"fmt"
	"github.com/memclutter/gocore/pkg/coreslices"
	"io/ioutil"
	"net/http"
)

func Get(options ...OptionFunc) error {
	return Run(append(options, WithMethod(http.MethodGet))...)
}

func Post(options ...OptionFunc) error {
	return Run(append(options, WithMethod(http.MethodPost))...)
}

func Put(options ...OptionFunc) error {
	return Run(append(options, WithMethod(http.MethodPut))...)
}

func Patch(options ...OptionFunc) error {
	return Run(append(options, WithMethod(http.MethodPatch))...)
}

func Delete(options ...OptionFunc) error {
	return Run(append(options, WithMethod(http.MethodDelete))...)
}

func Run(optionsOverride ...OptionFunc) (err error) {
	options := &Options{
		method:     http.MethodGet,
		extensions: make([]Extension, 0),
	}

	for _, optionOverride := range optionsOverride {
		if err := optionOverride(options); err != nil {
			return fmt.Errorf("options error: %v", err)
		}
	}

	client := http.DefaultClient
	for _, ext := range options.extensions {
		if client, err = ext.ClientOverride(client); err != nil {
			return fmt.Errorf("extension client override error: %v", err)
		}
	}

	req, err := http.NewRequest(options.method, options.url, nil)
	if err != nil {
		return fmt.Errorf("new request error: %v", err)
	}

	if options.headers != nil {
		for k, vv := range options.headers {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	if options.cookies != nil {
		for _, cookie := range options.cookies {
			req.AddCookie(cookie)
		}
	}

	for _, ext := range options.extensions {
		if req, err = ext.RequestOverride(req); err != nil {
			return fmt.Errorf("extension request override error: %v", err)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request error: %v", err)
	}
	defer res.Body.Close()

	// Read response
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read out error: %v", err)
	}

	// Check status codes
	if len(options.errStatusCodes) > 0 {
		if coreslices.IntIn(res.StatusCode, options.errStatusCodes) {
			return fmt.Errorf("http error status code %d: %s", res.StatusCode, out[:50])
		}
	}

	if len(options.okStatusCodes) > 0 {
		if !coreslices.IntIn(res.StatusCode, options.okStatusCodes) {
			return fmt.Errorf("http not ok status code %d: %s", res.StatusCode, out[:50])
		}
	}

	if options.out != nil {
		if options.outType == OutTypeBytes {
			if o, ok := options.out.(*[]byte); ok {
				*o = out
			} else {
				return fmt.Errorf("write out error: out is not *[]byte type")
			}
		} else if options.outType == OutTypeJson {
			if err := json.Unmarshal(out, options.out); err != nil {
				return fmt.Errorf("json unmarshal error: %v", err)
			}
		}
	}

	return nil
}
