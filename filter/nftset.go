package filter

import (
	"log"
	"net"
	"os/exec"
	"strings"
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

func parseNftsetList(buf []byte) (list []net.IP) {
	for _, line := range strings.Split(string(buf), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		if fields[0] != "add" {
			continue
		}
		ip := net.ParseIP(fields[2]).To4()
		if ip == nil {
			log.Printf("invalid IPv4 address: %q", fields[2])
			continue
		}
		list = append(list, ip)
	}
	return
}

func (b nftset) List() []net.IP {
	cmd := exec.Command("nft", "list", "set", "inet", "fw4", b.set)
	//cmd := exec.Command("ipset", "list", b.set, "-output", "save")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("command %q failed with %q", cmd, string(output))
		return nil
	}

	return parseNftsetList(output)
}
