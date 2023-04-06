package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// handle dynamic json keys
type Response struct {
	ZscloudNet map[string]map[string]city `json:"zscloud.net"`
}

type city []struct {
	Range string `json:"range"`
}

type Body []byte

type ResponseBody interface {
	io.ReadCloser
	Close() error
}

func (b Body) Close() error {
	return nil
}

// return slice "prefixes" containing ipv4 and ipv6 prefixes for each city for each continent
func ExtractIpPrefixes(b Body) ([]string, error) {
	var response Response

	// pointer automatically dereferenced (*r).Body
	if err := json.NewDecoder(b).Decode(&response); err != nil {
		return nil, err
	}

	var prefixes []string

	for _, continent := range response.ZscloudNet {
		for _, city := range continent {
			prefixes = append(prefixes, city[0].Range)
		}
	}

	return prefixes, nil
}

const zscaler = "https://api.config.zscaler.com/zscloud.net/cenr/json"

// break up for testing?
// create fixture for output testing
// push to git and pull
// add interface to convert based on argument -terraform -json (default to json)
func main() {

	res, err := http.Get(zscaler)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		if err != nil {
			log.Fatalln(err)
		}

		prefixes, err := ExtractIpPrefixes(res.Body)

		if err != nil {
			log.Fatalln(err)
		}

		// output consumable as input to terraform external data provider
		jsonOutput := fmt.Sprintf("{\"prefix_list\": \"%s\"}", strings.Join(prefixes, " "))
		fmt.Println(jsonOutput)
	}

}
