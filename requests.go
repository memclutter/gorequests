package gorequests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/memclutter/gocore/pkg/coreslices"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type requestsInstance struct {
	clientOverride    []ClientOverrideMiddleware
	requestOverride   []RequestOverrideMiddleware
	method            string
	url               string
	cookies           []*http.Cookie
	headers           http.Header
	data              []byte
	contentType       string
	form              url.Values
	json              interface{}
	responseOkCodes   []int
	responseFailCodes []int
	respRaw           *[]byte
	respJson          interface{}
}

func Trace(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodTrace).Url(url, args...)
}
func Connect(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodConnect).Url(url, args...)
}
func Head(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodHead).Url(url, args...)
}
func Options(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodOptions).Url(url, args...)
}
func Get(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodGet).Url(url, args...)
}
func Post(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodPost).Url(url, args...)
}
func Put(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodPut).Url(url, args...)
}
func Delete(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodDelete).Url(url, args...)
}
func Patch(url string, args ...any) RequestsInstance {
	return Requests().Method(http.MethodPatch).Url(url, args...)
}
func Requests() RequestsInstance { return new(requestsInstance) }

func (r *requestsInstance) Method(method string) RequestsInstance {
	r.method = method
	return r
}

func (r *requestsInstance) Url(url string, args ...any) RequestsInstance {
	r.url = fmt.Sprintf(url, args...)
	return r
}

func (r *requestsInstance) Cookies(cookies ...*http.Cookie) RequestsInstance {
	if r.cookies == nil {
		r.cookies = make([]*http.Cookie, 0)
	}
	for _, cookie := range cookies {
		r.cookies = append(r.cookies, cookie)
	}
	return r
}

func (r *requestsInstance) Header(key, value string) RequestsInstance {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Add(key, value)
	return r
}

func (r *requestsInstance) Data(data []byte, contentType ...string) RequestsInstance {
	r.data = data
	if len(contentType) > 0 {
		r.contentType = contentType[0]
	}
	return r
}

func (r *requestsInstance) Form(form url.Values) RequestsInstance {
	r.form = form
	return r
}

func (r *requestsInstance) Json(json interface{}) RequestsInstance {
	r.json = json
	return r
}

func (r *requestsInstance) ResponseCodeOk(codes ...int) RequestsInstance {
	r.responseOkCodes = codes
	return r
}

func (r *requestsInstance) ResponseCodeFail(codes ...int) RequestsInstance {
	r.responseFailCodes = codes
	return r
}

func (r *requestsInstance) ResponseRaw(responseRaw *[]byte) RequestsInstance {
	r.respRaw = responseRaw
	return r
}

func (r *requestsInstance) ResponseJson(respJson interface{}) RequestsInstance {
	r.respJson = respJson
	return r
}

func (r *requestsInstance) Use(middlewares ...interface{}) RequestsInstance {
	for _, middleware := range middlewares {
		if co, ok := middleware.(ClientOverrideMiddleware); ok {
			r.clientOverride = append(r.clientOverride, co)
		}
		if ro, ok := middleware.(RequestOverrideMiddleware); ok {
			r.requestOverride = append(r.requestOverride, ro)
		}
	}
	return r
}

func (r *requestsInstance) Exec() (err error) {
	var bodyReader io.Reader
	var contentType string
	if r.data != nil {
		bodyReader = bytes.NewReader(r.data)
		contentType = r.contentType
	}
	if r.form != nil {
		bodyReader = strings.NewReader(r.form.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else if r.json != nil {
		body, err := json.Marshal(r.json)
		if err != nil {
			return fmt.Errorf("request json body encode error: %v", err)
		}
		bodyReader = bytes.NewReader(body)
		contentType = "application/json"
	}

	// Create client instance and apply middleware
	c := &http.Client{}
	for _, co := range r.clientOverride {
		if c, err = co.ClientOverride(c); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return err
	}

	if len(contentType) != 0 {
		req.Header.Set("Content-Type", contentType)
	}
	if r.headers != nil {
		for k, vv := range r.headers {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	if r.cookies != nil {
		for _, cookie := range r.cookies {
			req.AddCookie(cookie)
		}
	}

	for _, ro := range r.requestOverride {
		if req, err = ro.RequestOverride(req); err != nil {
			return err
		}
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(r.responseFailCodes) > 0 {
		if coreslices.IntIn(res.StatusCode, r.responseFailCodes) {
			return fmt.Errorf(res.Status + string(body[:50]))
		}
	}
	if len(r.responseOkCodes) > 0 {
		if !coreslices.IntIn(res.StatusCode, r.responseOkCodes) {
			return fmt.Errorf(res.Status + string(body[:50]))
		}
	}
	if r.respRaw != nil {
		*(r.respRaw) = body
	}
	if r.respJson != nil {
		if err := json.Unmarshal(body, r.respJson); err != nil {
			return err
		}
	}
	return nil
}
