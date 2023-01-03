package gorequests

import "net/http"

type ClientOverrideMiddleware interface {
	ClientOverride(c *http.Client) (*http.Client, error)
}

type RequestOverrideMiddleware interface {
	RequestOverride(r *http.Request) (*http.Request, error)
}

type RequestsInstance interface {
	Use(middlewares ...interface{}) RequestsInstance
	Url(url string) RequestsInstance
	Method(method string) RequestsInstance
	Json(json interface{}) RequestsInstance
	Cookies(cookies ...*http.Cookie) RequestsInstance
	Header(key, value string) RequestsInstance
	ResponseCodeOk(codes ...int) RequestsInstance
	ResponseCodeFail(codes ...int) RequestsInstance
	ResponseRaw(responseRaw *[]byte) RequestsInstance
	ResponseJson(responseJson interface{}) RequestsInstance
	Exec() error
}
