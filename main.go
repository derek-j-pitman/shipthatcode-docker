package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PID Limit Enforcement
// pids.max per cgroup.
// Commands:
//   CGROUP <name>  -> create cgroup, no limit; print OK
//   PIDSMAX <name> <n>  -> set limit; print OK
//   FORK <name>  -> increment count; print 'OK <count>' or 'EAGAIN'
//   EXIT <name>  -> decrement count; print 'OK <count>'
//   STATUS <name>  -> print 'count=<n> max=<n_or_unlimited>'

type PidGroup struct {
	PidsMax, CurrPids int
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	groups := map[string]PidGroup{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "CGROUP":
			// TODO: create cgroup, no limit; print OK
			groups[parts[1]] = PidGroup{PidsMax: -1, CurrPids: 0}
			out = append(out, "OK")
		case "PIDSMAX":
			// TODO: set limit; print OK
			pidMax, _ := strconv.Atoi(parts[2])
			g := groups[parts[1]]
			g.PidsMax = pidMax
			groups[parts[1]] = g
			out = append(out, "OK")
		case "FORK":
			// TODO: increment count; print 'OK <count>' or 'EAGAIN'
			g := groups[parts[1]]
			if g.PidsMax >= 0 && g.CurrPids >= g.PidsMax {
				out = append(out, "EAGAIN")
			} else {
				g.CurrPids++
				out = append(out, fmt.Sprintf("OK %d", g.CurrPids))
			}
			groups[parts[1]] = g
		case "EXIT":
			// TODO: decrement count; print 'OK <count>'
			g := groups[parts[1]]
			g.CurrPids--
			if g.CurrPids < 0 {
				g.CurrPids = 0
			}
			groups[parts[1]] = g
			out = append(out, fmt.Sprintf("OK %d", g.CurrPids))
		case "STATUS":
			// TODO: print 'count=<n> max=<n_or_unlimited>'
			max := "unlimited"
			g := groups[parts[1]]
			if g.PidsMax >= 0 {
				max = strconv.Itoa(g.PidsMax)
			}
			out = append(out, fmt.Sprintf("count=%d max=%s", g.CurrPids, max))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
