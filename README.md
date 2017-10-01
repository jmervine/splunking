# Splunk Request

Low level lib to create an http Request object for connecting to Splunk.

### Usage

```go
import (
    "fmt"

    "github.com/jmervine/splunking"
)

func main() {
    req, err := splunking.Request("GET", "/api/path", nil)
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
