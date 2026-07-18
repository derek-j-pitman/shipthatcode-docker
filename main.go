package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// UTS Namespace Isolation
// Each UTS ns owns (hostname, domain). Defaults: ('localhost', '(none)').
// Commands:
//   NEWNS  -> new UTS ns; print id
//   SETHOST <ns> <name>  -> set hostname (max 64); OK or 'ERR too long'
//   SETDOMAIN <ns> <name>  -> set domain; OK
//   HOSTNAME <ns>  -> print hostname
//   DOMAIN <ns>  -> print domain or '(none)'
//   UNAME <ns>  -> print '<hostname> <domain>'

type UTS struct {
	Host, Domain string
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	var namespaces []UTS
	// TODO: declare your state structures here
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: new UTS ns; print id
			namespaces = append(namespaces, UTS{"localhost", "(none)"})
			out = append(out, strconv.Itoa(len(namespaces)))
		case "SETHOST":
			// TODO: set hostname (max 64); OK or 'ERR too long'
			nsId, _ := strconv.Atoi(parts[1])
			if len(parts[2]) > 64 {
				out = append(out, "ERR too long")
			} else {
				namespaces[nsId-1].Host = parts[2]
				out = append(out, "OK")
			}
		case "SETDOMAIN":
			// TODO: set domain; OK
			nsId, _ := strconv.Atoi(parts[1])
			namespaces[nsId-1].Domain = parts[2]
			out = append(out, "OK")
		case "HOSTNAME":
			// TODO: print hostname
			nsId, _ := strconv.Atoi(parts[1])
			out = append(out, namespaces[nsId-1].Host)
		case "DOMAIN":
			// TODO: print domain or '(none)'
			nsId, _ := strconv.Atoi(parts[1])
			out = append(out, namespaces[nsId-1].Domain)
		case "UNAME":
			// TODO: print '<hostname> <domain>'
			nsId, _ := strconv.Atoi(parts[1])
			out = append(out, fmt.Sprintf("%s %s", namespaces[nsId-1].Host, namespaces[nsId-1].Domain))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
