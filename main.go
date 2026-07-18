package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

// veth + Bridge + NAT Router
// Bridges with gateway IPs, netns with veth attachment, routing decision per packet.
// Commands:
//   BRIDGE <name> <gw_ip>  -> create bridge; OK
//   NETNS <name>  -> create netns; OK
//   VETH <bridge> <netns> <ns_ip>  -> create veth pair; OK
//   ROUTE <netns> <cidr>  -> add direct route; OK
//   SEND <netns> <dst_ip>  -> DIRECT, NAT, or NO ROUTE
//   SHOW-BRIDGE <name>  -> '<netns> <ns_ip>' sorted by netns

type Bridge struct {
	Gateway      net.IP
	AttachedNses map[string]net.IP
}

type NetNamespace struct {
	MyAddr  net.IP
	Gateway net.IP
	Routes  []*net.IPNet
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	bridges := map[string]Bridge{}
	namespaces := map[string]NetNamespace{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "BRIDGE":
			// TODO: create bridge; OK
			bridges[parts[1]] = Bridge{net.ParseIP(parts[2]), map[string]net.IP{}}
			out = append(out, "OK")
		case "NETNS":
			// TODO: create netns; OK
			namespaces[parts[1]] = NetNamespace{net.IPv4zero, net.IPv4zero, []*net.IPNet{}}
			out = append(out, "OK")
		case "VETH":
			// TODO: create veth pair; OK
			bridge := bridges[parts[1]]
			virtAddr := net.ParseIP(parts[3])
			bridge.AttachedNses[parts[2]] = virtAddr
			bridges[parts[1]] = bridge
			ns := namespaces[parts[2]]
			ns.MyAddr = virtAddr
			ns.Gateway = bridge.Gateway
			namespaces[parts[2]] = ns
			out = append(out, "OK")
		case "ROUTE":
			// TODO: add direct route; OK
			ns := namespaces[parts[1]]
			_, cidr, _ := net.ParseCIDR(parts[2])
			ns.Routes = append(ns.Routes, cidr)
			namespaces[parts[1]] = ns
			out = append(out, "OK")
		case "SEND":
			// TODO: DIRECT, NAT, or NO ROUTE
			dst := net.ParseIP(parts[2])
			ns := namespaces[parts[1]]
			rte := "NO ROUTE"
			for _, r := range ns.Routes {
				if r.Contains(dst) {
					rte = fmt.Sprintf("DIRECT from %s to %s", ns.MyAddr, dst)
				}
			}
			if rte == "NO ROUTE" && !ns.Gateway.Equal(net.IPv4zero) {
				rte = fmt.Sprintf("NAT from %s via %s to %s", ns.MyAddr, ns.Gateway, dst)
			}
			out = append(out, rte)
		case "SHOW-BRIDGE":
			// TODO: '<netns> <ns_ip>' sorted by netns
			bridge := bridges[parts[1]]
			var keys []string
			var veths []string
			for k := range bridge.AttachedNses {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				veths = append(veths, fmt.Sprintf("%s %s", k, namespaces[k].MyAddr))
			}
			out = append(out, strings.Join(veths, "\n"))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
