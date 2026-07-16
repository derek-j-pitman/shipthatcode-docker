package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Network Namespace Simulator
// Each netns has its own interfaces (starts with lo 127.0.0.1) and routes.
// Commands:
//   NEWNS  -> create new netns with lo 127.0.0.1; print id
//   IFACE-ADD <ns> <name> <ip>  -> add interface; print OK
//   IFACE-LIST <ns>  -> print '<name> <ip>' sorted by name
//   ROUTE-ADD <ns> <dst_cidr> <via_iface>  -> add route; print OK
//   ROUTE-LIST <ns>  -> print '<cidr> via <iface>' sorted by cidr

type Namespace struct {
	Ifaces map[string]string
	Routes map[string]string
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	var namespaces []Namespace
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: create new netns with lo 127.0.0.1; print id
			namespaces = append(namespaces, Namespace{Ifaces: map[string]string{"lo": "127.0.0.1"}, Routes: map[string]string{}})
			out = append(out, strconv.Itoa(len(namespaces)))
		case "IFACE-ADD":
			// TODO: add interface; print OK
			nsId, _ := strconv.Atoi(parts[1])
			name := parts[2]
			ip := parts[3]
			namespaces[nsId-1].Ifaces[name] = ip
			out = append(out, "OK")
		case "IFACE-LIST":
			// TODO: print '<name> <ip>' sorted by name
			nsId, _ := strconv.Atoi(parts[1])
			var ifaces []string
			var keys []string
			for k := range namespaces[nsId-1].Ifaces {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				ifaces = append(ifaces, fmt.Sprintf("%s %s", k, namespaces[nsId-1].Ifaces[k]))
			}
			out = append(out, strings.Join(ifaces, "\n"))
		case "ROUTE-ADD":
			// TODO: add route; print OK
			nsId, _ := strconv.Atoi(parts[1])
			cidr := parts[2]
			iface := parts[3]
			namespaces[nsId-1].Routes[cidr] = iface
			out = append(out, "OK")
		case "ROUTE-LIST":
			// TODO: print '<cidr> via <iface>' sorted by cidr
			nsId, _ := strconv.Atoi(parts[1])
			var routes []string
			var keys []string
			for k := range namespaces[nsId-1].Routes {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				routes = append(routes, fmt.Sprintf("%s via %s", k, namespaces[nsId-1].Routes[k]))
			}
			out = append(out, strings.Join(routes, "\n"))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
