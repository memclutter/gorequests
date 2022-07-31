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
