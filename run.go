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

func Run(optionsOverride ...OptionFunc) error {
	options := &Options{
		method: http.MethodGet,
	}

	for _, optionOverride := range optionsOverride {
		if err := optionOverride(options); err != nil {
			return fmt.Errorf("options error: %v", err)
		}
	}

	client := http.Client{}

	req, err := http.NewRequest(options.method, options.url, nil)
	if err != nil {
		return fmt.Errorf("new request error: %v", err)
	}

	if options.headers != nil {
		req.Header = options.headers
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
