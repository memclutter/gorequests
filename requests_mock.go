package gorequests

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type mockClientOverrideMiddleware struct {
	mock.Mock
}

func (m *mockClientOverrideMiddleware) ClientOverride(c *http.Client) (*http.Client, error) {
	args := m.Called(c)
	return args.Get(0).(*http.Client), args.Error(1)
}

type mockRequestOverrideMiddleware struct {
	mock.Mock
}

func (m *mockRequestOverrideMiddleware) RequestOverride(r *http.Request) (*http.Request, error) {
	args := m.Called(r)
	return args.Get(0).(*http.Request), args.Error(1)
}
