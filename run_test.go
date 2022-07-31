package gorequests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	go func() {
		if err := http.ListenAndServe(":8000", server{}); err != nil {
			panic(err.Error())
		}
	}()

	excepted := serverInfo{
		Method: http.MethodGet,
		Path:   "/path/",
	}
	out := serverInfo{}
	assert.Nil(t, Run(
		WithMethod(http.MethodGet),
		WithUrl("http://localhost:8000/path/"),
		WithOut(&out, OutTypeJson),
	), "The run should exec without error")
	assert.Equal(t, excepted, out, "The run should return correct json")
}

type serverInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type server struct {
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info, _ := json.Marshal(serverInfo{
		Method: r.Method,
		Path:   r.URL.Path,
	})
	w.Write(info)
}
