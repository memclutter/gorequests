package gorequests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func init() {
	// Run echo webserver
	go func() {
		if err := http.ListenAndServe(":8000", server{}); err != nil {
			panic(err.Error())
		}
	}()
}

func TestRun(t *testing.T) {

	excepted := serverInfo{
		Method: http.MethodGet,
		Path:   "/path/",
		Cookies: map[string]string{
			"id":   "1",
			"test": "test",
		},
		Headers: map[string][]string{
			"Accept":          []string{"application/json"},
			"Accept-Encoding": []string{"gzip"},
			"Authorization":   []string{"123-123-123"},
			"User-Agent":      []string{"Gorequests"},
		},
	}
	out := serverInfo{}
	assert.Nil(t, Run(
		WithMethod(http.MethodGet),
		WithUrl("http://localhost:8000/path/"),
		WithHeader("accept", "application/json"),
		WithHeader("accept-encoding", "gzip"),
		WithHeader("authorization", "123-123-123"),
		WithHeader("user-agent", "Gorequests"),
		WithCookies(&http.Cookie{Name: "id", Value: "1"}, &http.Cookie{Name: "test", Value: "test"}),
		WithOut(&out, OutTypeJson),
	), "The run should exec without error")
	assert.Equal(t, excepted, out, "The run should return correct json")
}

func TestGet(t *testing.T) {
	testMethod(t, Get, http.MethodGet)
}

func TestPost(t *testing.T) {
	testMethod(t, Post, http.MethodPost)
}

func TestPut(t *testing.T) {
	testMethod(t, Put, http.MethodPut)
}

func TestPatch(t *testing.T) {
	testMethod(t, Patch, http.MethodPatch)
}

func TestDelete(t *testing.T) {
	testMethod(t, Delete, http.MethodDelete)
}

func testMethod(t *testing.T, f func(...OptionFunc) error, method string) {
	out := serverInfo{}
	assert.Nil(t, f(
		WithUrl("http://localhost:8000/path/"),
		WithOut(&out, OutTypeJson),
	), "The run should exec without error")
	assert.Equalf(t, method, out.Method, "The %s should return correct method", method)
}

type serverInfo struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Cookies map[string]string   `json:"cookies"`
	Headers map[string][]string `json:"headers"`
}

type server struct {
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	headers := make(map[string][]string)
	for k, vv := range r.Header {
		if k == "Cookie" {
			continue
		}
		headers[k] = make([]string, 0)
		for _, v := range vv {
			headers[k] = append(headers[k], v)
		}
	}
	info, _ := json.Marshal(serverInfo{
		Method:  r.Method,
		Path:    r.URL.Path,
		Cookies: cookies,
		Headers: headers,
	})
	w.Write(info)
}
