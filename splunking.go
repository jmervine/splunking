package splunking

import (
	"fmt"
	"io"
	"net/http"

	"github.com/joeshaw/envdecode"
)

type SplunkRequest struct {
	Username   string `env:"SPLUNK_USERNAME,required",json:"username"`
	Password   string `env:"SPLUNK_PASSWORD,required",json:"password"`
	Host       string `env:"SPLUNK_HOST,required",json:"host"`
	Port       string `env:"SPLUNK_POST,default=8089",json:"port"`
	OutputMode string `env:"SPLUNK_OUTPUT_TYPE,default=json",json:"output_type"`
}

func Init() (SplunkRequest, error) {
	sr := SplunkRequest{}
	err := envdecode.Decode(&sr)

	return sr, err
}

func (sr *SplunkRequest) Get(endpoint string, body io.Reader) (resp *http.Response, err error) {
	req, err := sr.Request("GET", endpoint, body)
	if err != nil {
		return nil, err
	}

	return sr.Submit(req)
}

func (sr *SplunkRequest) Post(endpoint string, body io.Reader) (resp *http.Response, err error) {
	req, err := sr.Request("POST", endpoint, body)
	if err != nil {
		return nil, err
	}

	return sr.Submit(req)
}

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

func (sr *SplunkRequest) Submit(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func (sr *SplunkRequest) Endpoint(path string) string {
	return fmt.Sprintf("https://%s:%s%s", sr.Host, sr.Port, path)
}
