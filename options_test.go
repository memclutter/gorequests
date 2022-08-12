package gorequests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestWithClientTimeout(t *testing.T) {
	cases := []struct {
		clientTimeout time.Duration
	}{
		{clientTimeout: 10 * time.Second},
		{clientTimeout: 500 * time.Millisecond},
	}

	for _, c := range cases {
		options := &Options{}
		assert.NoError(t, WithClientTimeout(c.clientTimeout)(options))
		assert.Equal(t, c.clientTimeout, options.clientTimeout, "The options.clientTimeout should be equal")
	}
}

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
		assert.NoError(t, WithUrl(c.url, c.args...)(options))
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
		assert.NoError(t, WithMethod(c.method)(options))
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
		assert.NoError(t, WithHeader(c.key, c.value)(options))
		assert.Equal(t, c.value, options.headers.Get(c.key), "The options.header value should be set")
	}
}

func TestWithHeaders(t *testing.T) {
	cases := []struct {
		headers http.Header
	}{
		{headers: http.Header{
			"accept":       []string{"application/json"},
			"content-type": []string{"form/urlencoded"},
		}},
	}

	for _, c := range cases {
		options := &Options{}
		assert.NoError(t, WithHeaders(c.headers)(options))
		for k, vv := range c.headers {
			for _, v := range vv {
				assert.Equalf(t, v, options.headers.Get(k), "The options.header[%s] value should be set", k)
			}
		}
	}
}

func TestWithCookies(t *testing.T) {
	cases := []struct {
		cookies []*http.Cookie
	}{
		{
			cookies: []*http.Cookie{
				&http.Cookie{
					Name:     "test",
					Value:    "testval",
					Domain:   "example.com",
					HttpOnly: true,
					Path:     "/path/",
					Secure:   true,
				},
			},
		},
	}

	for _, c := range cases {
		options := &Options{}
		assert.NoError(t, WithCookies(c.cookies...)(options))
		for i, exceptedCookie := range c.cookies {
			found := false
			for _, cookie := range options.cookies {
				found = true
				if cookie.Name == exceptedCookie.Name {
					assert.Equalf(t, exceptedCookie.Value, cookie.Value, "The options.cookies[%d: %s].Value should be equal", i, exceptedCookie.Name)
					assert.Equalf(t, exceptedCookie.Domain, cookie.Domain, "The options.cookies[%d: %s].Domain should be equal", i, exceptedCookie.Name)
					assert.Equalf(t, exceptedCookie.HttpOnly, cookie.HttpOnly, "The options.cookies[%d: %s].HttpOnly should be equal", i, exceptedCookie.Name)
					assert.Equalf(t, exceptedCookie.Path, cookie.Path, "The options.cookies[%d: %s].Path should be equal", i, exceptedCookie.Name)
					assert.Equalf(t, exceptedCookie.Secure, cookie.Secure, "The options.cookies[%d: %s].Secure should be equal", i, exceptedCookie.Name)
				}
			}

			assert.Truef(t, found, "The options.cookies[%d: %s] should be found", i, exceptedCookie.Name)
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
		assert.NoError(t, WithOut(c.out, c.outType)(options))
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

func TestWithErrStatusCodes(t *testing.T) {
	cases := []struct {
		codes []int
	}{
		{
			codes: []int{http.StatusInternalServerError, http.StatusBadGateway},
		},
	}

	for _, c := range cases {
		options := &Options{}
		assert.NoError(t, WithErrStatusCodes(c.codes...)(options), "Must be run without errors")
		assert.Equal(t, c.codes, options.errStatusCodes, "The options.errStatusCodes should be equal")
	}
}

func TestWithOkStatusCodes(t *testing.T) {
	cases := []struct {
		codes []int
	}{
		{
			codes: []int{http.StatusOK, http.StatusCreated},
		},
	}

	for _, c := range cases {
		options := &Options{}
		assert.NoError(t, WithOkStatusCodes(c.codes...)(options), "Must be run without errors")
		assert.Equal(t, c.codes, options.okStatusCodes, "The options.okStatusCodes should be equal")
	}
}
