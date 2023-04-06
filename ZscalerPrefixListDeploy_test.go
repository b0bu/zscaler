package main

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"net"
	"strings"
	"testing"
)

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
