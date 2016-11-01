package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yawn/awhois"
)

type result struct {
	IP      net.IP
	Matches []string
}

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("usage: %s ip-address", os.Args[0])
	}

	ip := net.ParseIP(os.Args[1])

	if ip == nil {
		log.Fatalf("failed to parse IP address %s", os.Args[1])
	}

	matches, err := awhois.Check(ip)

	if err != nil {
		log.Fatal(err)
	}

	r := result{
		IP:      ip,
		Matches: make([]string, len(matches)),
	}

	for i, e := range matches {
		r.Matches[i] = e.String()
	}

	out, err := json.Marshal(r)

	if err != nil {
		log.Fatalf("failed to encode result to JSON: %s", err)
	}

	fmt.Fprintf(os.Stdout, string(out))

	if len(matches) == 0 {
		os.Exit(1)
	}

}
