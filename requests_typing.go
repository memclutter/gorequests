package gorequests

import (
	"net/http"
	"net/url"
)

type ClientOverrideMiddleware interface {
	ClientOverride(c *http.Client) (*http.Client, error)
}

type RequestOverrideMiddleware interface {
	RequestOverride(r *http.Request) (*http.Request, error)
}

type RequestsShort func(url string, args ...any) RequestsInstance

type RequestsInstance interface {
	Use(middlewares ...interface{}) RequestsInstance
	Url(url string, args ...interface{}) RequestsInstance
	Method(method string) RequestsInstance
	Data(data []byte, contentType ...string) RequestsInstance
	Form(form url.Values) RequestsInstance
	Json(json interface{}) RequestsInstance
	Cookies(cookies ...*http.Cookie) RequestsInstance
	Header(key, value string) RequestsInstance
	ResponseCodeOk(codes ...int) RequestsInstance
	ResponseCodeFail(codes ...int) RequestsInstance
	ResponseRaw(responseRaw *[]byte) RequestsInstance
	ResponseJson(responseJson interface{}) RequestsInstance
	Exec() error
}
