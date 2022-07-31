package gorequests

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"time"
)

type Extension interface {
	ClientOverride(c *http.Client) (*http.Client, error)
	RequestOverride(req *http.Request) (*http.Request, error)
}

// Retry extension for example

type RetryExtension struct {
	RetryMax     int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

func (e RetryExtension) ClientOverride(c *http.Client) (*http.Client, error) {
	rc := retryablehttp.NewClient()
	rc.RetryMax = e.RetryMax
	rc.RetryWaitMin = e.RetryWaitMin
	rc.RetryWaitMax = e.RetryWaitMax
	return rc.StandardClient(), nil
}

func (e RetryExtension) RequestOverride(req *http.Request) (*http.Request, error) { return req, nil }
