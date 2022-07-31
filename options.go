package gorequests

import (
	"fmt"
	"net/http"
)

type OutType int

const (
	OutTypeBytes = iota
	OutTypeJson
	//OutTypeXml @TODO...
)

type Options struct {
	method  string
	url     string
	headers http.Header
	out     interface{}
	outType OutType
}

type OptionFunc func(options *Options) error

func WithMethod(method string) OptionFunc {
	return func(options *Options) error {
		options.method = method
		return nil
	}
}

func WithUrl(url string, args ...interface{}) OptionFunc {
	return func(options *Options) error {
		if len(args) > 0 {
			url = fmt.Sprintf(url, args...)
		}
		options.url = url
		return nil
	}
}

func WithHeaders(headers http.Header) OptionFunc {
	return func(options *Options) error {
		if options.headers == nil {
			options.headers = headers
		} else {
			for k, vv := range headers {
				for _, v := range vv {
					options.headers.Add(k, v)
				}
			}
		}
		return nil
	}
}

func WithHeader(key, value string) OptionFunc {
	return func(options *Options) error {
		if options.headers == nil {
			options.headers = http.Header{}
		}
		options.headers.Add(key, value)
		return nil
	}
}

func WithOut(out interface{}, outType OutType) OptionFunc {
	return func(options *Options) error {
		options.out = out
		options.outType = outType
		return nil
	}
}
