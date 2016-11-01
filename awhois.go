package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type match struct {
	network *net.IPNet
	region  string
	service string
}

func (m *match) String() string {
	return fmt.Sprintf("%s (%s %s)", m.network, m.service, m.region)
}

type response struct {
	IPv6Prefixes []map[string]string `json:"ipv6_prefixes"`
	Prefixes     []map[string]string
}

func main() {

	const (
		reg = "region"
		srv = "service"
		v4p = "ip_prefix"
		v6p = "ipv6_prefix"
	)

	if len(os.Args) != 2 {
		log.Fatalf("usage: %s :ip-address", os.Args[0])
	}

	var (
		ip      = net.ParseIP(os.Args[1])
		key     string
		matches []*match
		targets []map[string]string
		v4      = ip.To4() != nil
	)

	r, err := fetch()

	if err != nil {
		log.Fatalf("failed to retrieve IP prefixes: %s", err)
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
			log.Fatalf("unexpected data at %s", e)
		}

		_, network, err := net.ParseCIDR(prefix)

		if err != nil {
			log.Fatalf("failed to parse prefix %q: %s", prefix, err)
		}

		if network.Contains(ip) {

			matches = append(matches, &match{
				network: network,
				region:  region,
				service: service,
			})

		}

	}

	if len(matches) == 0 {
		log.Printf("%s is not part of AWS", ip)
		os.Exit(1)
	}

	log.Printf("%s is part of AWS: %s", ip, matches)

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
		return nil, err
	}

	return response, nil

}
