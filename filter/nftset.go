package filter

import (
	"log"
	"net"
	"os/exec"
	"github.com/buger/jsonparser"
)

type nftset struct {
	set string
}

func NewNftset(set string) *nftset {
	cmd := exec.Command("nft", "add", "set", "inet", "fw4", set, "{ type ipv4_addr ; }")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
	}
	return &nftset{set}
}

func (b nftset) Block(ip net.IP) {
	cmd := exec.Command("nft", "add", "element", "inet", "fw4", b.set, "{ " + ip.String() + "  }")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
	}
}

func (b nftset) Unblock(ip net.IP) {
	cmd := exec.Command("nft", "delete", "element", "inet", "fw4", b.set, "{ " + ip.String() + " }")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
	}
}

func parseNftsetList(buf []byte) (list []net.IP) {

	if parsed, _, _, err := jsonparser.Get(buf, "nftables"); err == nil {
		buf = parsed
	} else {
		return nil
	}

	jsonparser.ArrayEach(buf, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if setdata, _, _, err := jsonparser.Get(value, "set"); err == nil {
			if elements, _, _, err := jsonparser.Get(setdata, "elem"); err == nil {
				jsonparser.ArrayEach(elements, func(value2 []byte, dataType2 jsonparser.ValueType, offset2 int, err2 error) {
					if dataType2 == jsonparser.String {
						if str, err := jsonparser.ParseString(value2); err == nil {
							if ip := net.ParseIP(str).To4(); ip != nil {
								list = append(list, ip)
							} else {
								log.Printf("invalid IPv4 address: %q", str)
							}
						}
					}
				})
			}
		}
	})

	return
}

func (b nftset) List() []net.IP {
        cmd := exec.Command("nft", "--json", "list", "set", "inet", "fw4", b.set)
        output, err := cmd.CombinedOutput()
        if err != nil {
                log.Printf("command %q failed with %q", cmd, string(output))
                return nil
        }

	return parseNftsetList(output)
}
