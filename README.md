# Go Soap

package to help with SOAP integrations (client)

### Install

```bash
go get github.com/noho-digital/gosoap
```

### Configuration Options

gosoap can be configured when creating the SOAP client by using `SoapClientWithConfig()`
instead of the regular `SoapClient()` method. This method takes an additional `gosoap.Config{}`
struct that has the following configuration options:

- Dump (`boolean`): this makes gosoap dump the raw HTTP request and response in the log (useful for debugging)
- Logger (`DumpLogger`): Logger takes any type that implements the [`gosoap.DumpLogger`](https://github.com/Siteminds/gosoap/blob/master/soap.go#L29) interface. This is useful for wrapping Dump logs into your own logging solution (e.g. zap, logrus, zerolog, etc.)
- PrefixOperation (`boolean`): by default gosoap generates an 'operation' root element in the SOAP Body, containing a `xmlns` namespace attribute. When setting `PrefixOperation = true`, the `xmlns` attribute is not added, allowing you to manually prefix the operation name in the `gosoap.Call()` with your own namespace prefix. You can add your own namespace prefix using `gosoap.SetCustomEnvelope()` (see example below).
- DisableRoot (`boolean`): when set to true, this makes gosoap skip generating the operation root element entirely, allowing full control of the SOAP Body.
- Endpoint (`string`): by default the first location of the first service port found in the WSDL is used. By defining the endpoint url in the configuration, this configured endpoint will be used instead of the one from the WSDL.

### Examples

#### Basic use

```go
package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/noho-digital/gosoap"
)

// GetIPLocationResponse will hold the Soap response
type GetIPLocationResponse struct {
	GetIPLocationResult string `xml:"GetIpLocationResult"`
}

// GetIPLocationResult will
type GetIPLocationResult struct {
	XMLName xml.Name `xml:"GeoIP"`
	Country string   `xml:"Country"`
	State   string   `xml:"State"`
}

var (
	r GetIPLocationResponse
)

func main() {
	httpClient := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	soap, err := gosoap.SoapClient("http://wsgeoip.lavasoft.com/ipservice.asmx?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	// Use gosoap.ArrayParams to support fixed position params
	params := gosoap.Params{
		"sIp": "8.8.8.8",
	}

	res, err := soap.Call("GetIpLocation", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	res.Unmarshal(&r)

	// GetIpLocationResult will be a string. We need to parse it to XML
	result := GetIPLocationResult{}
	err = xml.Unmarshal([]byte(r.GetIPLocationResult), &result)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	if result.Country != "US" {
		log.Fatalf("error: %+v", r)
	}

	log.Println("Country: ", result.Country)
	log.Println("State: ", result.State)
}
```

#### Set Custom Envelope Attributes

```go
package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/noho-digital/gosoap"
)

// GetIPLocationResponse will hold the Soap response
type GetIPLocationResponse struct {
	GetIPLocationResult string `xml:"GetIpLocationResult"`
}

// GetIPLocationResult will
type GetIPLocationResult struct {
	XMLName xml.Name `xml:"GeoIP"`
	Country string   `xml:"Country"`
	State   string   `xml:"State"`
}

var (
	r GetIPLocationResponse
)

func main() {
	httpClient := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	// set custom envelope
    gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:tem": "http://tempuri.org/",
    })

	soap, err := gosoap.SoapClient("http://wsgeoip.lavasoft.com/ipservice.asmx?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	// Use gosoap.ArrayParams to support fixed position params
	params := gosoap.Params{
		"sIp": "8.8.8.8",
	}

	res, err := soap.Call("GetIpLocation", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	res.Unmarshal(&r)

	// GetIpLocationResult will be a string. We need to parse it to XML
	result := GetIPLocationResult{}
	err = xml.Unmarshal([]byte(r.GetIPLocationResult), &result)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	if result.Country != "US" {
		log.Fatalf("error: %+v", r)
	}

	log.Println("Country: ", result.Country)
	log.Println("State: ", result.State)
}
```

#### Set Header params

```go
	soap.HeaderParams = gosoap.SliceParams{
		xml.StartElement{
			Name: xml.Name{
				Space: "auth",
				Local: "Login",
			},
		},
		"user",
		xml.EndElement{
			Name: xml.Name{
				Space: "auth",
				Local: "Login",
			},
		},
		xml.StartElement{
			Name: xml.Name{
				Space: "auth",
				Local: "Password",
			},
		},
		"P@ssw0rd",
		xml.EndElement{
			Name: xml.Name{
				Space: "auth",
				Local: "Password",
			},
		},
	}
```
