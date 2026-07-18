package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Memory cgroup Simulator
// Memory cgroup with hard limit. OOM if alloc would exceed limit.
// Commands:
//   CGROUP <name>  -> create cgroup, no limit; print OK
//   LIMIT <name> <bytes>  -> set hard limit; print OK
//   ALLOC <name> <pid> <bytes>  -> alloc; print OK or 'OOM <pid>' (don't add bytes on OOM)
//   FREE <name> <pid> <bytes>  -> free; print OK
//   STATUS <name>  -> print 'usage=<n> limit=<n_or_unlimited>'

type CGroup struct {
	Limit  int
	Usage  int
	PTable map[int]int
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
			// TODO: create cgroup, no limit; print OK
			groups[parts[1]] = CGroup{-1, 0, map[int]int{}}
			out = append(out, "OK")
		case "LIMIT":
			// TODO: set hard limit; print OK
			newLimit, _ := strconv.Atoi(parts[2])
			g := groups[parts[1]]
			g.Limit = newLimit
			groups[parts[1]] = g
			out = append(out, "OK")
		case "ALLOC":
			// TODO: alloc; print OK or 'OOM <pid>' (don't add bytes on OOM)
			newPid, _ := strconv.Atoi(parts[2])
			pSize, _ := strconv.Atoi(parts[3])
			if groups[parts[1]].Limit > 0 && (pSize+groups[parts[1]].Usage) > groups[parts[1]].Limit {
				out = append(out, fmt.Sprintf("OOM %d", newPid))
			} else {
				g := groups[parts[1]]
				g.PTable[newPid] += pSize
				g.Usage += pSize
				groups[parts[1]] = g
				out = append(out, "OK")
			}
		case "FREE":
			// TODO: free; print OK
			g := groups[parts[1]]
			pid, _ := strconv.Atoi(parts[2])
			toFree, _ := strconv.Atoi(parts[3])
			if g.PTable[pid] <= toFree {
				g.Usage -= g.PTable[pid]
				delete(g.PTable, pid)
			} else {
				g.PTable[pid] -= toFree
				g.Usage -= toFree
			}
			groups[parts[1]] = g
			out = append(out, "OK")
		case "STATUS":
			// TODO: print 'usage=<n> limit=<n_or_unlimited>'
			limit := "unlimited"
			if groups[parts[1]].Limit > 0 {
				limit = strconv.Itoa(groups[parts[1]].Limit)
			}
			out = append(out, fmt.Sprintf("usage=%s limit=%s", strconv.Itoa(groups[parts[1]].Usage), limit))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
