package splunking

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/joeshaw/envdecode"
)

type SplunkRequest struct {
	Username   string `env:"SPLUNK_USERNAME,required" json:"username"`
	Password   string `env:"SPLUNK_PASSWORD,required" json:"password"`
	Host       string `env:"SPLUNK_HOST,required" json:"host"`
	Port       string `env:"SPLUNK_POST,default=8089" json:"port"`
	Proto      string `env:"SPLUNK_PROTO,default=https" json:"proto"`
	OutputMode string `env:"SPLUNK_OUTPUT_TYPE,default=json" json:"output_type"`
}

// InitURL allows for initializing with a base url. Expected format examples
// include:
//    InitURL("user:pass@host")
//    InitURL("https://user:pass@host:port?output_mode=mode")
//
// Default port is '8089' and default output_mode is 'json'. 'https' will be
// prepended if a protocol isn't passed.
func InitURL(str string) (sr SplunkRequest, err error) {
	// Check for proto, it's required to parse username and password correctly
	if !strings.HasPrefix(str, "https://") && !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "//") {
		str = "https://" + str // default to https
	}

	u, err := url.Parse(str)
	if err != nil {
		return
	}

	// Ensure that proto is always http or https.
	if u.Scheme == "http" {
		sr.Proto = "http"
	} else {
		sr.Proto = "https"
	}

	if u.User == nil {
		err = errors.New("Username is required")
		return
	}

	sr.Username = u.User.Username()
	if sr.Username == "" {
		err = errors.New("Username is required")
		return
	}

	var ok bool
	sr.Password, ok = u.User.Password()
	if !ok || sr.Password == "" {
		err = errors.New("Password is required")
		return
	}

	split := strings.Split(u.Host, ":")
	sr.Host = split[0]

	if sr.Host == "" {
		err = errors.New("Host is required")
		return
	}

	sr.Port = "8089" // default port
	if len(split) > 1 {
		sr.Port = split[1]
	}

	sr.OutputMode = u.Query().Get("output_mode")
	if sr.OutputMode == "" {
		sr.OutputMode = "json"
	}

	return
}

// Init loads configuration from the environment.
//
//     SPLUNK_USERNAME=username
//     SPLUNK_PASSWORD=password
//     SPLUNK_HOST=splunk.example.com
//     SPLUNK_PORT=8089        // default
//     SPLUNK_PROTO=https      // default
//     SPLUNK_OUTPUT_TYPE=json // default
func Init() (SplunkRequest, error) {
	sr := SplunkRequest{}
	err := envdecode.Decode(&sr)

	// Ensure that proto is always http or https.
	if sr.Proto != "http" && sr.Proto != "https" {
		sr.Proto = "https"
	}

	return sr, err
}

func (sr *SplunkRequest) simpleRequest(method, endpoint string, body io.Reader) (resp *http.Response, err error) {
	req, err := sr.Request(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	return sr.Submit(req)
}

// Get wraps a simple http GET, setting up the request and calling submitting it.
func (sr *SplunkRequest) Get(endpoint string, body io.Reader) (resp *http.Response, err error) {
	return sr.simpleRequest("GET", endpoint, body)
}

// Post wraps an http POST, setting up the request and calling submitting it.
func (sr *SplunkRequest) Post(endpoint string, body io.Reader) (resp *http.Response, err error) {
	return sr.simpleRequest("POST", endpoint, body)
}

// Delete wraps an http DELETE, setting up the request and calling submitting it.
func (sr *SplunkRequest) Delete(endpoint string, body io.Reader) (resp *http.Response, err error) {
	return sr.simpleRequest("DELETE", endpoint, body)
}

// Request initializes a base request, building the endpoint, setting up auth,
// adding the correct headers and including the output_type.
func (sr *SplunkRequest) Request(method, endpoint string, body io.Reader) (req *http.Request, err error) {
	endpoint = sr.Endpoint(endpoint)

	req, err = http.NewRequest(method, endpoint, body)
	if err != nil {
		return
	}

	req.SetBasicAuth(sr.Username, sr.Password)
	req.Header.Add("Accept", "application/json")

	params := req.URL.Query()
	if params.Get("output_mode") == "" {
		params.Add("output_mode", sr.OutputMode)
	}

	req.URL.RawQuery = params.Encode()

	return
}

// Submit is the follow up to a Request call, executing the request properly.
func (sr *SplunkRequest) Submit(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

// Endpoint generates a base URL for http interations with Splunk.
func (sr *SplunkRequest) Endpoint(path string) string {
	return fmt.Sprintf("%s://%s:%s%s", sr.Proto, sr.Host, sr.Port, path)
}
