package awhois

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type response struct {
	IPv6Prefixes []map[string]string `json:"ipv6_prefixes"`
	Prefixes     []map[string]string
}

// Check tests if a given IP has AWS matches. No matches are not an error but an
// empty slice.
func Check(ip net.IP) ([]*Match, error) {

	const (
		reg = "region"
		srv = "service"
		v4p = "ip_prefix"
		v6p = "ipv6_prefix"
	)

	var (
		key     string
		matches []*Match
		targets []map[string]string
		v4      = ip.To4() != nil
	)

	r, err := fetch()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve IP prefixes: %s", err)
	}

	if v4 {
		key = v4p
		targets = r.Prefixes
	} else {
		key = v6p
		targets = r.IPv6Prefixes
	}

	for _, e := range targets {

		var (
			ok      bool
			prefix  string
			region  = e[reg]
			service = e[srv]
		)

		if prefix, ok = e[key]; !ok {
			return nil, fmt.Errorf("unexpected data at %s", e)
		}

		_, network, err := net.ParseCIDR(prefix)

		if err != nil {
			return nil, fmt.Errorf("failed to parse prefix %q: %s", prefix, err)
		}

		if network.Contains(ip) {

			matches = append(matches, &Match{
				Network: network,
				Region:  region,
				Service: service,
			})

		}

	}

	return matches, nil

}

func fetch() (*response, error) {

	const url = "https://ip-ranges.amazonaws.com/ip-ranges.json"

	var response = new(response)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %s", err)
	}

	return response, nil

}
