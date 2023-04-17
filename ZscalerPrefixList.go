package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const zscalerNetApiEndpoint = "https://api.config.zscaler.com/zscloud.net/cenr/json"

// handle dynamic json keys
type Response struct {
	ZsCloudNet map[string]map[string]City `json:"zscloud.net"`
}

type City []struct {
	Range string `json:"range"`
}

// filter cities by allow list
func allowedRegion(s string) bool {
	regions := []string{"Amsterdam", "London", "Paris", "Manchester"}
	for _, region := range regions {
		if strings.Contains(s, region) {
			return true
		}
	}
	return false
}

// return slice "prefixes" containing ipv4 and ipv6 prefixes for each city for each continent
func ExtractIpPrefixes(r Response) []string {
	var prefixes []string
	for _, continent := range r.ZsCloudNet {
		for city, ranges := range continent {
			if allowedRegion(city) {
				prefixes = append(prefixes, ranges[0].Range)
			}
		}
	}
	return prefixes
}

func ToStructE(r *http.Response) Response {
	var response Response

	// pointer automatically dereferenced (*res).Body
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		log.Fatalln(err)
	}
	return response
}

func HttpGet(endpoint string) (*http.Response, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}
	return nil, fmt.Errorf("expected 200 response got %d", resp.StatusCode)
}

func main() {

	resp, err := HttpGet(zscalerNetApiEndpoint)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	prefixes := ExtractIpPrefixes(ToStructE(resp))

	// output consumable as input to terraform external data provider
	jsonOutput := fmt.Sprintf("{\"prefix_list\": \"%s\"}", strings.Join(prefixes, " "))
	fmt.Println(jsonOutput)

}
