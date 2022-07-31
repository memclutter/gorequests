package gorequests

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"math/rand"
	"net/http"
	"net/url"
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
	rc.Logger = nil
	return rc.StandardClient(), nil
}

func (e RetryExtension) RequestOverride(req *http.Request) (*http.Request, error) { return req, nil }

// Proxies extensions

type ProxiesExtension struct{ Proxies []string }

func (e ProxiesExtension) ClientOverride(c *http.Client) (*http.Client, error) {
	if len(e.Proxies) == 0 {
		return nil, fmt.Errorf("nil proxies")
	}
	proxy := e.Proxies[rand.Intn(len(e.Proxies))]
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("error parse proxy url %s: %v", proxy, err)
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}, nil
}

func (e ProxiesExtension) RequestOverride(req *http.Request) (*http.Request, error) { return req, nil }
