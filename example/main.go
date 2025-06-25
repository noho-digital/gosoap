package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
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
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current working directory: %s", err)
	}

	// set custom envelope
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:sch":     "http://lavasoft.com/",
	})

	soap, err := gosoap.SoapClientWithConfig(
		fmt.Sprintf("file://%s/../testdata/ipservice.wsdl", pwd),
		client(),
		&gosoap.Config{
			Dump:            true,
			PrefixOperation: true,
			Endpoint:        "http://localhost:8080",
		},
	)

	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	// Use gosoap.ArrayParams to support fixed position params
	params := gosoap.Params{
		"sch:sIp": "8.8.8.8",
	}

	res, err := soap.Call("sch:GetIpLocation", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	if err = res.Unmarshal(&r); err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

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

// setup http client with sane timeouts and transport config
func client() *http.Client {
	return &http.Client{
		Transport: transport(),
		Timeout:   5 * time.Second,
	}
}

// setup tcp transport with sane timeouts and tls config
func transport() *http.Transport {
	return &http.Transport{
		TLSClientConfig:     tlsConfig(),
		TLSHandshakeTimeout: 2 * time.Second,
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
	}
}

// disable strict TLS certificate checking..
func tlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}
