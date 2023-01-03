package gorequests

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type RequestSuite struct {
	suite.Suite
}

func (suite *RequestSuite) BeforeTest(suiteName, testName string) {
	httpmock.Activate()
}

func (suite *RequestSuite) AfterTest(suiteName, testName string) {
	httpmock.Reset()
}

func (suite *RequestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (suite *RequestSuite) TestPerformGetOk() {
	// Test data
	method := http.MethodGet
	reqUrl := "https://localhost:9000/"
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, httpmock.NewStringResponder(http.StatusOK, `{"ip": "127.0.0.1"}`))

	// Run test target
	err := Request().Method(method).Url(reqUrl).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equalf(suite.T(), 1, httpStats[callKey], "should be call once")
}

func (suite *RequestSuite) TestPerformPost() {
	// Test data
	method := http.MethodPost
	reqUrl := "https://localhost:9000/projects"
	body := map[string]string{"name": "Test", "description": "Awesome test project"}
	callKey := method + " " + reqUrl

	// Mocking http calls
	httpmock.RegisterResponder(method, reqUrl, func(request *http.Request) (*http.Response, error) {
		exceptedBody, _ := json.Marshal(body)
		body, err := ioutil.ReadAll(request.Body)
		if err == nil && strings.Contains(string(body), string(exceptedBody)) {
			return httpmock.NewStringResponder(http.StatusOK, `{"id": 1}`)(request)
		}
		return httpmock.ConnectionFailure(request)
	})

	// Run test target
	err := Request().Method(method).Url(reqUrl).Json(body).Exec()

	// Prepare assert stats
	httpStats := httpmock.GetCallCountInfo()

	// Assertions
	assert.NoError(suite.T(), err, "should be run without error")
	assert.Contains(suite.T(), httpStats, callKey, "should be send http request")
	assert.Equal(suite.T(), 1, httpStats[callKey], "should be call once")
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}
