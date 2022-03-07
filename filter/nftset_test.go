package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseNsetList(t *testing.T) {
	input := []byte(`
{
  "nftables": [
    {
      "metainfo": {
        "version": "1.0.2",
        "release_name": "Lester Gooch",
        "json_schema_version": 1
      }
    },
    {
      "set": {
        "family": "inet",
        "name": "myset",
        "table": "fw4",
        "type": "ipv4_addr",
        "handle": 501,
        "elem": [
          "1.1.1.1",
          "2.2.2.2",
          "3.3.3.3"
        ]
      }
    }
  ]
}
`)
	expected := []string{
		"1.1.1.1",
		"2.2.2.2",
		"3.3.3.3",
	}

	list := parseNftsetList(input)
	output := make([]string, len(list))
	for i, l := range list {
		output[i] = l.String()
	}
	assert.Equal(t, expected, output)
}
