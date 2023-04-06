package main

import (
	"bytes"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"net"
	"strings"
	"testing"
)

var testData = []Body(`"
{
	"zscaler.net":{
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
		  ],
		  "city : Amsterdam II":[
			 {
				"range":"147.161.172.0/23",
				"vpn":"",
				"gre":"",
				"hostname":"",
				"latitude":"52",
				"longitude":"5"
			 },
			 {
				"range":"165.225.240.0/23",
				"vpn":"ams2-2-vpn.zscaler.net",
				"gre":"165.225.240.12",
				"hostname":"ams2-2.sme.zscaler.net",
				"latitude":"52",
				"longitude":"5"
			 }
		  ]
	   },
	   "continent : EMEA":{
		  "city : Abu Dhabi I":[
			 {
				"range":"147.161.175.0/23",
				"vpn":"",
				"gre":"",
				"hostname":"",
				"latitude":"24.453884",
				"longitude":"54.3773438"
			 }
		  ],
		  "city : Amsterdam II":[
			 {
				"range":"147.161.173.0/23",
				"vpn":"",
				"gre":"",
				"hostname":"",
				"latitude":"52",
				"longitude":"5"
			 },
			 {
				"range":"165.225.250.0/23",
				"vpn":"ams2-2-vpn.zscaler.net",
				"gre":"165.225.240.12",
				"hostname":"ams2-2.sme.zscaler.net",
				"latitude":"52",
				"longitude":"5"
			 }
		  ]
	   }
	}
 }
"`)

// Ensure api response is parsed correctly
func TestExtractIpPrefixes(t *testing.T) {
	t.Parallel()

	prefixes, err := ExtractIpPrefixes(testData)
	if err != nil {
		t.Errorf("Could no unmarshall json.")
	}

	// marshal test data to struct

	// assert actual with got
	fmt.Println(prefixes)
}

// Ensure that zscaler ips are filtered ipv4 only
func TestActualTerraformOutputIsIp4(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "fixtures/",
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
