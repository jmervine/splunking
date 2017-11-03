package splunking

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestInitURL(t *testing.T) {
	assert := assert.New(t)

	expect := SplunkRequest{"foo", "bar", "example.com", "9999", "http", "xml"}
	got, err := InitURL("http://foo:bar@example.com:9999?output_mode=xml")

	assert.Nil(err)
	assert.Equal(expect, got)

	expect = SplunkRequest{"foo", "bar", "example.com", "8089", "https", "json"}

	// url.Parse will error without a proto, so this tests proto prepending
	// in addition to default port and output_mode
	got, err = InitURL("foo:bar@example.com")

	assert.Nil(err)
	assert.Equal(expect, got)

	_, err = InitURL("https://example.com")
	assert.Equal(errors.New("Username is required"), err)

	_, err = InitURL("https://:bar@example.com")
	assert.Equal(errors.New("Username is required"), err)

	_, err = InitURL("https://foo:@example.com")
	assert.Equal(errors.New("Password is required"), err)

	_, err = InitURL("https://foo:bar@")
	assert.Equal(errors.New("Host is required"), err)
}

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

func TestEndpoint(t *testing.T) {
	assert := assert.New(t)

	sr1 := SplunkRequest{"user1", "pass1", "host1.com", "9999", "http", ""}
	sr2 := SplunkRequest{"user1", "pass1", "host1.com", "9999", "https", ""}

	assert.Equal("http://host1.com:9999/foo/bar", sr1.Endpoint("/foo/bar"))
	assert.Equal("https://host1.com:9999/foo/bar", sr2.Endpoint("/foo/bar"))
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
		Proto:      "https",
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
