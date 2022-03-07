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
	//cmd := exec.Command("ipset", "create", set, "hash:ip")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
	}
	return &nftset{set}
}

func (b nftset) Block(ip net.IP) {
	cmd := exec.Command("nft", "add", "element", "inet", "fw4", b.set, "{ " + ip.String() + "  }")
	//cmd := exec.Command("ipset", "add", b.set, ip.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
	}
}

func (b nftset) Unblock(ip net.IP) {
	cmd := exec.Command("nft", "delete", "element", "inet", "fw4", b.set, "{ " + ip.String() + " }")
	//cmd := exec.Command("ipset", "del", b.set, ip.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
	}
}
/*
func parseNftsetList(buf []byte) (list []net.IP) {
	for _, line := range strings.Split(string(buf), "\n") {
		ip := net.ParseIP(line).To4()
		if ip == nil {
			log.Printf("invalid IPv4 address: %q", line)
			continue
		}
		list = append(list, ip)
	}
	return
}
*/
func (b nftset) List() []net.IP {
	cmd := exec.Command("nft", "--json", "list", "set", "inet", "fw4", b.set)
	//cmd := exec.Command("ipset", "list", b.set, "-output", "save")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
		return nil
	}

	if parsed, _, _, err := jsonparser.Get(output, "nftables"); err == nil {
		output = parsed
	} else {
		return nil
	}

	var elemlist []net.IP
	success := false

	jsonparser.ArrayEach(output, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if setdata, _, _, err := jsonparser.Get(value, "set"); err == nil {
			if elements, _, _, err := jsonparser.Get(setdata, "elem"); err == nil {
				var elemlist []net.IP
				jsonparser.ArrayEach(elements, func(value2 []byte, dataType2 jsonparser.ValueType, offset2 int, err2 error) {
					if dataType2 == jsonparser.String {
						if str, err := jsonparser.ParseString(value2); err == nil {
							if ip := net.ParseIP(str).To4(); ip != nil {
								elemlist = append(elemlist, ip)
							} else {
								log.Printf("invalid IPv4 address: %q", str)
							}
						}
					}
				success = true
				})
			}
		}
	})

	if success {
		return elemlist
	}

	return nil
}
