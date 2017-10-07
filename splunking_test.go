package splunking

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestRequest(t *testing.T) {
	assert := assert.New(t)

	sr := defaultRequest()
	r, e := sr.Request("GET", "/api/path", nil)

	assert.Nil(e)
	assert.Equal(r.Header.Get("Accept"), "application/json")
	assert.True(strings.HasPrefix(r.Header.Get("Authorization"), "Basic "))
	assert.Equal(r.URL.String(), "https://splunk.example.com:8089/api/path?output_mode=json")
}

func TestSubmit(t *testing.T) {
	defer mockRequest("GET", "https://splunk.example.com:8089/api/path?output_mode=json")()

	sr := defaultRequest()
	req, _ := sr.Request("GET", "/api/path", nil)
	resp, err := sr.Submit(req)

	assertResponse(t, resp, err, "at=GET")
}

// TODO: Also testing params handling, should be broken out in to it's own
// test.
func TestGet(t *testing.T) {
	defer mockRequest("GET", "https://splunk.example.com:8089/api/path?foo=bar&output_mode=somethingelse")()

	sr := defaultRequest()
	resp, err := sr.Get("/api/path?foo=bar&output_mode=somethingelse", nil)

	assertResponse(t, resp, err, "at=GET")
}

func TestPost(t *testing.T) {
	defer mockRequest("POST", "https://splunk.example.com:8089/api/path?output_mode=json")()

	sr := defaultRequest()
	resp, err := sr.Post("/api/path", nil)

	assertResponse(t, resp, err, "at=POST")
}

func TestDelete(t *testing.T) {
	defer mockRequest("DELETE", "https://splunk.example.com:8089/api/path?output_mode=json")()

	sr := defaultRequest()

	resp, err := sr.Delete("/api/path", nil)

	assertResponse(t, resp, err, "at=DELETE")
}

func mockRequest(method, url string) func() {
	httpmock.Activate()
	httpmock.RegisterResponder(method, url, httpmock.NewStringResponder(200, "at="+method))

	return httpmock.DeactivateAndReset
}

func defaultRequest() *SplunkRequest {
	return &SplunkRequest{
		Username:   "username",
		Password:   "password",
		Host:       "splunk.example.com",
		Port:       "8089",
		OutputMode: "json",
	}
}

func assertResponse(t *testing.T, r *http.Response, e error, expected string) {
	if assert.Nil(t, e) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, string(body), expected)
	}
}
