package splunking

import (
	//"fmt"
	"strings"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestRequest(t *testing.T) {
	sr := defaultRequest()
	r, e := sr.Request("GET", "/api/path", nil)

	if e != nil {
		t.Error("Expected nil, got,", e)
	}

	expected := r.Header.Get("Accept")
	if expected != "application/json" {
		t.Error("Expected application/json got,", expected)
	}

	expected = r.Header.Get("Authorization")
	if !strings.HasPrefix(expected, "Basic ") {
		t.Error("Expected Basic Authication got,", expected)
	}

	expected = "https://splunk.example.com:8089/api/path?output_mode=json"
	if r.URL.String() != expected {
		t.Error("Expected", expected, " got,", expected)
	}
}

func TestSubmit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://splunk.example.com:8089/api/path?output_mode=json",
		httpmock.NewStringResponder(200, `{}`))

	sr := defaultRequest()
	req, _ := sr.Request("GET", "/api/path", nil)

	_, e := sr.Submit(req)
	if e != nil {
		t.Error("Expected nil, got,", e)
	}
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
