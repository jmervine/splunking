# Splunk Request

[![Build Status](https://travis-ci.org/jmervine/splunking.svg?branch=master)](https://travis-ci.org/jmervine/splunking) [![GoDoc](https://godoc.org/github.com/jmervine/splunking?status.svg)](https://godoc.org/github.com/jmervine/splunking)


Low level lib to create an http Request object for connecting to Splunk.

### Usage

```go
import (
    "fmt"

    "github.com/jmervine/splunking"
)

func main() {
    // Load configs from environment
    //  SPLUNK_USERNAME=username
    //  SPLUNK_PASSWORD=password
    //  SPLUNK_HOST=splunk.example.com
    //  SPLUNK_PORT=8089        // default
    //  SPLUNK_PROTO=https      // default
    //  SPLUNK_OUTPUT_TYPE=json // default
    client, err := splunking.Init()
    if err != nil {
        panic(err)
    }

    // Or load from URL - same defaults as above apply
    client, err = splunking.InitURL("https://username:password@splunk.example.com:8089?output_type=json")
    if err != nil {
        panic(err)
    }

    // Submit a request.
    req, err := client.Request("GET", "/api/path", nil)
    if err != nil {
        panic(err)
    }

    // ... do stuff with request

    resp, err := splunking.Submit(req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", resp)
}
```
