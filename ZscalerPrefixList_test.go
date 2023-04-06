package main

import (
	"bytes"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var testData = `
{
	"zscloud.net":{
	   "continent : EMEA":{
		  "city : Abu Dhabi I":[
			 {
				"range":"147.161.174.0/23",
				"vpn":"",
				"gre":"",
				"hostname":"",
				"latitude":"24.453884",
				"longitude":"54.3773438"
			 }
		  ]
	   }
	}
 }`

// Ensure api response is parsed correctly
func TestExtractIpPrefixes(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(testData))),
	}

	got := ExtractIpPrefixes(ToStructE(resp))
	wants := []string{"147.161.174.0/23"}

	assert.ElementsMatch(t, got, wants)
}

// Ensure http json response to struct
func TestToStructE(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(testData))),
	}

	responseLiteral := Response{map[string]map[string]City{
		"continent : EMEA": {
			"city : Abu Dhabi I": City{
				{"147.161.174.0/23"},
			},
		},
	}}

	response := ToStructE(resp)
	assert.True(t, reflect.DeepEqual(response, responseLiteral))
}

// Ensure that zscaler ips are filtered ipv4 only
func TestActualTerraformOutputIsIp4(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "tests/fixtures/",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.OutputList(t, terraformOptions, "zscaler_ipv4_ranges")

	for _, ip := range output {
		ipWithoutPrefixLength := strings.Split(ip, "/")[0]
		ip := net.ParseIP(ipWithoutPrefixLength)
		if ip.To4() == nil {
			t.Errorf("%v is not a valid ipv4 address.\n", ip)
		}
	}
}
