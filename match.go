package awhois

import (
	"fmt"
	"net"
)

// Match describes a matching network
type Match struct {
	Network *net.IPNet
	Region  string
	Service string
}

func (m *Match) String() string {
	return fmt.Sprintf("%s (%s %s)", m.Network, m.Service, m.Region)
}
