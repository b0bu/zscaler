package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const zscalerNetApi = "https://api.config.zscaler.com/zscloud.net/cenr/json"

// handle dynamic json keys
type Response struct {
	ZsCloudNet map[string]map[string]City `json:"zscloud.net"`
}

type City []struct {
	Range string `json:"range"`
}

// return slice "prefixes" containing ipv4 and ipv6 prefixes for each city for each continent
func ExtractIpPrefixes(r Response) []string {
	var prefixes []string

	for _, continent := range r.ZsCloudNet {
		for _, city := range continent {
			prefixes = append(prefixes, city[0].Range)
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

// break up for testing?
// create fixture for output testing
// push to git and pull
// add interface to convert based on argument -terraform -json (default to json)
func main() {

	resp, err := http.Get(zscalerNetApi)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%#v\n", resp.Body)
		prefixes := ExtractIpPrefixes(ToStructE(resp))

		// output consumable as input to terraform external data provider
		jsonOutput := fmt.Sprintf("{\"prefix_list\": \"%s\"}", strings.Join(prefixes, " "))
		fmt.Println(jsonOutput)
	}

}
