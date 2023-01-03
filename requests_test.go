package gorequests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

type RequestsSuite struct{ suite.Suite }

func (suite *RequestsSuite) BeforeTest(suiteName, testName string) { httpmock.Activate() }
func (suite *RequestsSuite) AfterTest(suiteName, testName string)  { httpmock.Reset() }
func (suite *RequestsSuite) TearDownSuite()                        { httpmock.DeactivateAndReset() }
func TestRequestSuite(t *testing.T)                                { suite.Run(t, new(RequestsSuite)) }

func (suite *RequestsSuite) TestShortcuts() {
	tests := []struct {
		method string
		target func(url string) RequestsInstance
	}{
		{method: http.MethodGet, target: Get},
		{method: http.MethodPost, target: Post},
		{method: http.MethodPut, target: Put},
		{method: http.MethodPatch, target: Patch},
		{method: http.MethodHead, target: Head},
		{method: http.MethodOptions, target: Options},
		{method: http.MethodDelete, target: Delete},
		{method: http.MethodTrace, target: Trace},
		{method: http.MethodConnect, target: Connect},
	}

	for _, test := range tests {
		suite.Run(test.method, func() {
			// Test data
			reqUrl := "http://localhost"
			callKey := test.method + " " + reqUrl

			// Mocking http calls
			httpmock.RegisterResponder(test.method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string("")))

			// Run test target
			err := test.target(reqUrl).Exec()

			// Prepare assert stats
			httpStats := httpmock.GetCallCountInfo()

			// Assertions
			assert.NoError(suite.T(), err, "should be run without error")
			assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
			assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
		})
	}
}

func (suite *RequestsSuite) TestCookies() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	reqCookies := []*http.Cookie{
		{
			Name:  "token",
			Value: "test",
		},
	}
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, func(request *http.Request) (*http.Response, error) {
		actualCookies := request.Cookies()
		if len(reqCookies) == len(actualCookies) {
			return httpmock.NewStringResponder(http.StatusOK, string(""))(request)
		}
		return httpmock.ConnectionFailure(request)
	})

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Cookies(reqCookies...).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestsSuite) TestHeader() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	reqHeader := "X-Test"
	reqHeaderValue := "test"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, func(request *http.Request) (*http.Response, error) {
		if request.Header.Get(reqHeader) == reqHeaderValue {
			return httpmock.NewStringResponder(http.StatusOK, string(""))(request)
		}
		return httpmock.ConnectionFailure(request)
	})

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Header(reqHeader, reqHeaderValue).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestsSuite) TestJson() {
	// Test data
	method := http.MethodPost
	reqUrl := "http://localhost/books"
	callKey := method + " " + reqUrl
	reqJson := map[string]string{"title": "A book", "description": "About a ..."}
	reqJsonBytes, _ := json.Marshal(reqJson)

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, func(request *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(request.Body)
		if request.Header.Get("Content-Type") == "application/json" && bytes.Contains(body, reqJsonBytes) {
			return httpmock.NewStringResponder(http.StatusOK, string(""))(request)
		}
		return httpmock.ConnectionFailure(request)
	})

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Json(reqJson).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestsSuite) TestResponseOkCodes() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusForbidden, ""))

	// Run test target
	err := Requests().Url(reqUrl).Method(method).ResponseCodeOk(http.StatusOK).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.Error(suite.T(), err, "should be run with error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestsSuite) TestResponseFailCodes() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusForbidden, ""))

	// Run test target
	err := Requests().Url(reqUrl).Method(method).ResponseCodeFail(http.StatusForbidden).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.Error(suite.T(), err, "should be run with error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestsSuite) TestResponseRaw() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl
	responseRaw := []byte(`SUCCESS`)

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string(responseRaw)))

	// Run test target
	actualResponseRaw := make([]byte, 0)
	err := Requests().Url(reqUrl).Method(method).ResponseRaw(&actualResponseRaw).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
	assert.Equal(suite.T(), responseRaw, actualResponseRaw, "should be equal response")
}

func (suite *RequestsSuite) TestResponseJson() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl
	responseJson := map[string]bool{"success": true}
	responseRaw, _ := json.Marshal(responseJson)

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string(responseRaw)))

	// Run test target
	actualResponseJson := make(map[string]bool)
	err := Requests().Url(reqUrl).Method(method).ResponseJson(&actualResponseJson).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
	assert.Equal(suite.T(), responseJson, actualResponseJson, "should be equal response")
}

func (suite *RequestsSuite) TestResponseJsonErr() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl
	responseRaw := []byte(`Server error`)

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusInternalServerError, string(responseRaw)))

	// Run test target
	actualResponseJson := make(map[string]bool)
	err := Requests().Url(reqUrl).Method(method).ResponseJson(&actualResponseJson).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.Error(suite.T(), err, "should be run with error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestsSuite) TestUseRequestOverrideMiddleware() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string("")))

	// Expected request
	expectedRequest, _ := http.NewRequest(method, reqUrl, nil)

	// Mocking request override middleware
	mockMiddleware := new(mockRequestOverrideMiddleware)
	mockMiddleware.On("RequestOverride", expectedRequest).Return(expectedRequest, nil)

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Use(mockMiddleware).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
	mockMiddleware.AssertExpectations(suite.T())
}

func (suite *RequestsSuite) TestUseRequestOverrideMiddlewareErr() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string("")))

	// Expected request
	expectedRequest, _ := http.NewRequest(method, reqUrl, nil)

	// Mocking request override middleware
	mockMiddleware := new(mockRequestOverrideMiddleware)
	mockMiddleware.On("RequestOverride", expectedRequest).Return(expectedRequest, fmt.Errorf("test"))

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Use(mockMiddleware).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.Error(suite.T(), err, "should be run with error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 0, httpStats[callKey], "should be call once")
	mockMiddleware.AssertExpectations(suite.T())
}

func (suite *RequestsSuite) TestUseClientOverrideMiddleware() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string("")))

	// Expected client
	expectedClient := &http.Client{}

	// Mocking request override middleware
	mockMiddleware := new(mockClientOverrideMiddleware)
	mockMiddleware.On("ClientOverride", expectedClient).Return(expectedClient, nil)

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Use(mockMiddleware).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
	mockMiddleware.AssertExpectations(suite.T())
}

func (suite *RequestsSuite) TestUseClientOverrideMiddlewareErr() {
	// Test data
	method := http.MethodGet
	reqUrl := "http://localhost"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, string("")))

	// Expected client
	expectedClient := &http.Client{}

	// Mocking request override middleware
	mockMiddleware := new(mockClientOverrideMiddleware)
	mockMiddleware.On("ClientOverride", expectedClient).Return(expectedClient, fmt.Errorf("test"))

	// Run test target
	err := Requests().Url(reqUrl).Method(method).Use(mockMiddleware).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.Error(suite.T(), err, "should be run with error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 0, httpStats[callKey], "should be call once")
	mockMiddleware.AssertExpectations(suite.T())
}
