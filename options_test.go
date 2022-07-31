package gorequests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithUrl(t *testing.T) {
	cases := []struct {
		url  string
		args []interface{}
	}{
		{
			url:  "https://example.com/%s/path/%d/id/",
			args: []interface{}{"some", 10},
		},
	}

	for _, c := range cases {
		options := &Options{}
		if err := WithUrl(c.url, c.args...)(options); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fmt.Sprintf(c.url, c.args...), options.url, "The options.url should be equal")
	}
}

func TestWithMethod(t *testing.T) {
	cases := []struct{ method string }{
		{method: http.MethodGet},
		{method: http.MethodPost},
		{method: http.MethodDelete},
		{method: http.MethodPut},
		{method: http.MethodOptions},
		{method: http.MethodHead},
	}

	for _, c := range cases {
		options := &Options{}
		if err := WithMethod(c.method)(options); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, c.method, options.method, "The options.method should be equal")
	}
}

func TestWithHeader(t *testing.T) {
	cases := []struct {
		key   string
		value string
	}{
		{key: "accept", value: "application/json"},
		{key: "content-type", value: "text/plain"},
	}

	for _, c := range cases {
		options := &Options{}
		if err := WithHeader(c.key, c.value)(options); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, c.value, options.headers.Get(c.key), "The options.header value should be set")
	}
}

func TestWithHeaders(t *testing.T) {
	cases := []struct {
		headers http.Header
	}{
		{headers: http.Header{
			"Accept":       []string{"application/json"},
			"Content-Type": []string{"form/urlencoded"},
		}},
	}

	for _, c := range cases {
		options := &Options{}
		if err := WithHeaders(c.headers)(options); err != nil {
			t.Fatal(err)
		}
		for k, vv := range c.headers {
			for _, v := range vv {
				assert.Equalf(t, v, options.headers.Get(k), "The options.header[%s] value should be set", k)
			}
		}
	}
}

func TestWithOut(t *testing.T) {
	cases := []struct {
		outType OutType
		out     interface{}
	}{
		{
			outType: OutTypeBytes,
			out:     &[]byte{},
		},
		{
			outType: OutTypeJson,
			out:     &map[string]string{},
		},
	}

	for _, c := range cases {
		options := &Options{}
		if err := WithOut(c.out, c.outType)(options); err != nil {
			t.Fatal(err)
		}
		if c.outType == OutTypeBytes {
			if out, ok := options.out.(*[]byte); ok {
				*out = []byte("test")
				assert.Equal(t, []byte("test"), *(c.out.(*[]byte)), "The options.out value should be equal[bytes]")
			} else {
				t.Fatal("not bytes in out")
			}
		} else if c.outType == OutTypeJson {
			data := map[string]string{"test": "test"}
			newData := c.out.(*map[string]string)
			*newData = data
			assert.Equal(t, data, *(c.out.(*map[string]string)), "The options.out value should be equal[json]")
		}
	}
}
