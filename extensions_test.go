package gorequests

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRetryExtension_ClientOverride(t *testing.T) {
	e := RetryExtension{
		RetryMax:     10,
		RetryWaitMin: 5 * time.Second,
		RetryWaitMax: 10 * time.Second,
	}
	c, err := e.ClientOverride(http.DefaultClient)
	assert.NoError(t, err, "Client override should run without errors")
	assert.IsType(t, c.Transport, &retryablehttp.RoundTripper{}, "Client should override transport")
	rc := c.Transport.(*retryablehttp.RoundTripper).Client
	assert.NotNil(t, rc, "Retryable client should not nil")
	assert.Equal(t, e.RetryMax, rc.RetryMax, "Retry max should be equal")
	assert.Equal(t, e.RetryWaitMin, rc.RetryWaitMin, "Retry wait min should be equal")
	assert.Equal(t, e.RetryWaitMax, rc.RetryWaitMax, "Retry wait max should be equal")
}

func TestProxiesExtension_ClientOverride(t *testing.T) {
	e := ProxiesExtension{
		Proxies: []string{"socks5://user:pass@host:8000"},
	}
	c, err := e.ClientOverride(http.DefaultClient)
	assert.NoError(t, err, "Client override should run without error")
	assert.IsType(t, c.Transport, &http.Transport{}, "Client should override transport")
	tr := c.Transport.(*http.Transport)
	proxy, err := tr.Proxy(&http.Request{})
	assert.NoError(t, err, "Proxy should run without error")
	pw, _ := proxy.User.Password()
	assert.Equal(t, proxy.Scheme, "socks5", "Protocol should be equal")
	assert.Equal(t, proxy.User.Username(), "user", "User should be equal")
	assert.Equal(t, pw, "pass", "Password should be equal")
	assert.Equal(t, proxy.Host, "host:8000", "Host should be equal")
}
