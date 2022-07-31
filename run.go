package gorequests

import (
	"encoding/json"
	"fmt"
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

	client := &http.Client{}
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
