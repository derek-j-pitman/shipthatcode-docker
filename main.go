package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CPU Quota Tracker
// Quota in microseconds per 100ms period. Throttle when exceeded.
// Commands:
//   CGROUP <name>  -> default quota=100000; print OK
//   QUOTA <name> <us>  -> set quota; print OK
//   RUN <name> <us>  -> consume; print OK or 'THROTTLE <name>' (clamp used to quota)
//   TICK  -> reset all used to 0; print OK
//   STATUS <name>  -> print 'used=<n> quota=<n> throttled=<true|false>'

type CGroup struct {
	Quota, Used int
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	groups := map[string]CGroup{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "CGROUP":
			// TODO: default quota=100000; print OK
			groups[parts[1]] = CGroup{Quota: 100000, Used: 0}
			out = append(out, "OK")
		case "QUOTA":
			// TODO: set quota; print OK
			g := groups[parts[1]]
			newQuota, _ := strconv.Atoi(parts[2])
			g.Quota = newQuota
			groups[parts[1]] = g
			out = append(out, "OK")
		case "RUN":
			// TODO: consume; print OK or 'THROTTLE <name>' (clamp used to quota)
			g := groups[parts[1]]
			newUs, _ := strconv.Atoi(parts[2])
			g.Used = min(g.Quota, g.Used+newUs)
			if g.Used == g.Quota {
				out = append(out, fmt.Sprintf("THROTTLE %s", parts[1]))
			} else {
				out = append(out, "OK")
			}
			groups[parts[1]] = g
		case "TICK":
			// TODO: reset all used to 0; print OK
			for k := range groups {
				g := groups[k]
				g.Used = 0
				groups[k] = g
			}
			out = append(out, "OK")
		case "STATUS":
			// TODO: print 'used=<n> quota=<n> throttled=<true|false>'
			g := groups[parts[1]]
			out = append(out, fmt.Sprintf("used=%d quota=%d throttled=%t", g.Used, g.Quota, g.Used == g.Quota))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
