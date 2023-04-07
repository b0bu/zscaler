package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"reflect"
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

func TestHttpGet(t *testing.T) {
	resp, _ := HttpGet("http://checkip.amazonaws.com")
	defer resp.Body.Close()
	assert.True(t, resp.StatusCode == 200)
}
