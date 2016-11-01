package awhois

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {

	assert := assert.New(t)

	var tt = []struct {
		in  string
		out string
		err *string
	}{
		{"54.76.43.209", "54.76.0.0/15 (AMAZON eu-west-1) 54.76.0.0/15 (EC2 eu-west-1)", nil},
		{"2a05:d050:8080:0:0:0:0:0", "2a05:d050:8080::/44 (AMAZON eu-west-1)", nil},
	}

	for _, e := range tt {

		out, err := Check(net.ParseIP(e.in))

		if e.err != nil {
			assert.EqualError(err, *e.err)
		} else {

			var m []string

			for _, e := range out {
				m = append(m, e.String())
			}

			assert.Equal(e.out, strings.Join(m, " "))
			assert.NoError(err)
		}

	}

}
